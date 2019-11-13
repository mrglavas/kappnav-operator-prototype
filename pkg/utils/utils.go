package utils

import (
	kappnavv1 "github.com/kappnav/operator/pkg/apis/kappnav/v1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extensionsv1beta1 "k8s.io/api/extensions/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// KappnavExtension extends the reconciler to manage
// additional resources and override default configuration.
type KappnavExtension interface {
	ApplyAdditionalDefaults(instance *kappnavv1.Kappnav)
	ReconcileAdditionalResources(request reconcile.Request, instance *kappnavv1.Kappnav) (reconcile.Result, error)
}

const (
	// APIContainerName ...
	APIContainerName string = "kappnav-api"
	// UIContainerName ...
	UIContainerName string = "kappnav-ui"
	// ControllerContainerName ...
	ControllerContainerName string = "kappnav-controller"
	// ServiceAccountNameSuffix ...
	ServiceAccountNameSuffix string = "sa"
)

const (
	// DefaultAPIRepository ...
	DefaultAPIRepository kappnavv1.Repository = "kappnav/apis"
	// DefaultUIRepository ...
	DefaultUIRepository kappnavv1.Repository = "kappnav/ui"
	// DefaultControllerRepository ...
	DefaultControllerRepository kappnavv1.Repository = "kappnav/controller"
	// DefaultTag ...
	DefaultTag kappnavv1.Tag = "0.1.0"
)

const (
	// CPUConstraintDefault ...
	CPUConstraintDefault string = "500m"
	// MemoryConstraintDefault ...
	MemoryConstraintDefault string = "512Mi"
)

const (
	// KubeEnvDefault ...
	KubeEnvDefault string = "okd"
)

const (
	// UIVolumeName ...
	UIVolumeName string = "ui-service-tls"
	// UIVolumeMountPath ...
	UIVolumeMountPath string = "/etc/tls/private"
)

// GetLabels ...
func GetLabels(instance *kappnavv1.Kappnav,
	existingLabels map[string]string, component *metav1.ObjectMeta) map[string]string {
	labels := map[string]string{
		"app.kubernetes.io/name":       instance.Name,
		"app.kubernetes.io/instance":   instance.Name,
		"app.kubernetes.io/managed-by": "kappnav-operator",
	}
	if component != nil && len(component.Name) > 0 {
		labels["app.kubernetes.io/component"] = component.GetName()
	}
	// Allow app.kubernetes.io/name to be overriden by the CR.
	// See: https://github.com/appsody/appsody-operator/issues/179
	for key, value := range instance.Labels {
		if key != "app.kubernetes.io/instance" &&
			key != "app.kubernetes.io/component" &&
			key != "app.kubernetes.io/managed-by" {
			labels[key] = value
		}
	}
	if existingLabels == nil {
		return labels
	}
	// Add labels to the existing map.
	for key, value := range labels {
		existingLabels[key] = value
	}
	return existingLabels
}

// CustomizeServiceAccount ...
func CustomizeServiceAccount(sa *corev1.ServiceAccount, instance *kappnavv1.Kappnav) {
	sa.Labels = GetLabels(instance, sa.Labels, &sa.ObjectMeta)
	imagePullSecrets := make([]corev1.LocalObjectReference, 1)
	imagePullSecrets[0] = corev1.LocalObjectReference{
		Name: "sa-" + sa.GetNamespace(),
	}
	pullSecrets := instance.Spec.Image.PullSecrets
	if pullSecrets != nil && len(pullSecrets) != 0 {
		for _, secretName := range pullSecrets {
			imagePullSecrets = append(imagePullSecrets, corev1.LocalObjectReference{
				Name: secretName,
			})
		}
	}
	sa.ImagePullSecrets = imagePullSecrets
}

// CustomizeClusterRoleBinding ...
func CustomizeClusterRoleBinding(crb *rbacv1.ClusterRoleBinding,
	sa *corev1.ServiceAccount, instance *kappnavv1.Kappnav) {
	crb.Labels = GetLabels(instance, crb.Labels, &crb.ObjectMeta)
	crb.Subjects = []rbacv1.Subject{
		{
			Kind:      "ServiceAccount",
			Name:      sa.GetName(),
			Namespace: sa.GetNamespace(),
		},
	}
	crb.RoleRef = rbacv1.RoleRef{
		Kind:     "ClusterRole",
		Name:     "cluster-admin",
		APIGroup: "rbac.authorization.k8s.io",
	}
}

// CustomizeConfigMap ...
func CustomizeConfigMap(configMap *corev1.ConfigMap, instance *kappnavv1.Kappnav) {
	configMap.Labels = GetLabels(instance, configMap.Labels, &configMap.ObjectMeta)
}

// CustomizeSecret ...
func CustomizeSecret(secret *corev1.Secret, instance *kappnavv1.Kappnav) {
	secret.Labels = GetLabels(instance, secret.Labels, &secret.ObjectMeta)
}

// CustomizeService ...
func CustomizeService(service *corev1.Service, instance *kappnavv1.Kappnav, annotations map[string]string) {
	service.Labels = GetLabels(instance, service.Labels, &service.ObjectMeta)
	if service.Annotations == nil {
		service.Annotations = annotations
	} else {
		// Add annotations to the existing map.
		for key, value := range annotations {
			service.Annotations[key] = value
		}
	}
}

// CustomizeUIServiceSpec ...
func CustomizeUIServiceSpec(serviceSpec *corev1.ServiceSpec, instance *kappnavv1.Kappnav) {
	isMinikube := IsMinikubeEnv(instance.Spec.Env.KubeEnv)
	oldType := serviceSpec.Type
	if isMinikube {
		serviceSpec.Type = corev1.ServiceTypeNodePort
	} else {
		serviceSpec.Type = ""
	}
	if oldType != serviceSpec.Type {
		serviceSpec.Ports = nil
	}
	if serviceSpec.Ports == nil || len(serviceSpec.Ports) == 0 {
		if isMinikube {
			serviceSpec.Ports = []corev1.ServicePort{
				{
					Port:       3000,
					TargetPort: intstr.FromInt(3000),
					Protocol:   corev1.ProtocolTCP,
					Name:       "https",
				},
			}
		} else {
			serviceSpec.Ports = []corev1.ServicePort{
				{
					Name:       "proxy",
					Port:       443,
					TargetPort: intstr.FromInt(8443),
				},
			}
		}
	}
	serviceSpec.Selector = map[string]string{
		"app.kubernetes.io/component": instance.GetName() + "-ui",
	}
}

// CustomizeIngress ...
func CustomizeIngress(ingress *extensionsv1beta1.Ingress, instance *kappnavv1.Kappnav) {
	ingress.Labels = GetLabels(instance, ingress.Labels, &ingress.ObjectMeta)
}

// CustomizeUIIngressSpec ...
func CustomizeUIIngressSpec(ingressSpec *extensionsv1beta1.IngressSpec,
	uiService *corev1.Service, instance *kappnavv1.Kappnav) {
	if ingressSpec.Rules == nil || len(ingressSpec.Rules) == 0 {
		ingressSpec.Rules = []extensionsv1beta1.IngressRule{
			{
				IngressRuleValue: extensionsv1beta1.IngressRuleValue{
					HTTP: &extensionsv1beta1.HTTPIngressRuleValue{
						Paths: []extensionsv1beta1.HTTPIngressPath{
							{
								Path: "/kappnav-ui",
								Backend: extensionsv1beta1.IngressBackend{
									ServiceName: uiService.GetName(),
									ServicePort: intstr.FromInt(3000),
								},
							},
							{
								Path: "/kappnav",
								Backend: extensionsv1beta1.IngressBackend{
									ServiceName: uiService.GetName(),
									ServicePort: intstr.FromInt(3000),
								},
							},
						},
					},
				},
			},
		}
	}
}

// CustomizeRoute ...
func CustomizeRoute(route *routev1.Route, instance *kappnavv1.Kappnav) {
	route.Labels = GetLabels(instance, route.Labels, &route.ObjectMeta)
}

// CustomizeUIRouteSpec ...
func CustomizeUIRouteSpec(routeSpec *routev1.RouteSpec,
	routeName *metav1.ObjectMeta, instance *kappnavv1.Kappnav) {
	if routeSpec.TLS == nil {
		routeSpec.TLS = &routev1.TLSConfig{}
	}
	routeSpec.TLS.Termination = routev1.TLSTerminationReencrypt
	routeSpec.To.Kind = "Service"
	routeSpec.To.Name = routeName.GetName()
}

// CustomizeDeployment ...
func CustomizeDeployment(deploy *appsv1.Deployment, instance *kappnavv1.Kappnav) {
	deploy.Labels = GetLabels(instance, deploy.Labels, &deploy.ObjectMeta)
	// Ensure that there's at least one replica
	if deploy.Spec.Replicas == nil || *deploy.Spec.Replicas < 1 {
		one := int32(1)
		deploy.Spec.Replicas = &one
	}
	deploy.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			"app.kubernetes.io/component": deploy.GetName(),
		},
	}
}

// CustomizePodSpec ...
func CustomizePodSpec(pts *corev1.PodTemplateSpec, parentComponent *metav1.ObjectMeta,
	containers []corev1.Container, volumes []corev1.Volume, instance *kappnavv1.Kappnav) {
	pts.Labels = GetLabels(instance, pts.Labels, parentComponent)
	pts.Spec.Containers = containers
	pts.Spec.RestartPolicy = corev1.RestartPolicyAlways
	pts.Spec.ServiceAccountName = instance.GetName() + "-" + ServiceAccountNameSuffix
	pts.Spec.Volumes = volumes
	setPodSecurity(pts)
}

// CustomizeKappnavConfigMap ...
func CustomizeKappnavConfigMap(kappnavConfig *corev1.ConfigMap, instance *kappnavv1.Kappnav) {
	// Initialize the config map or restore values if they have been deleted.
	if kappnavConfig.Data == nil {
		kappnavConfig.Data = make(map[string]string)
	}
	value, _ := kappnavConfig.Data["status-color-mapping"]
	if len(value) == 0 {
		kappnavConfig.Data["status-color-mapping"] =
			"{ \"values\": { \"Normal\": \"GREEN\", \"Warning\": \"YELLOW\", \"Problem\": \"RED\", \"Unknown\": \"GREY\"}," +
				"\"colors\": { \"GREEN\":  \"#5aa700\", \"YELLOW\": \"#B4B017\", \"RED\": \"#A74343\", \"GREY\" : \"#808080\"} }"
	}
	value, _ = kappnavConfig.Data["app-status-precedence"]
	if len(value) == 0 {
		kappnavConfig.Data["app-status-precedence"] = "[ \"Problem\", \"Warning\", \"Unknown\", \"Normal\" ]"
	}
	value, _ = kappnavConfig.Data["status-unknown"]
	if len(value) == 0 {
		kappnavConfig.Data["status-unknown"] = "Unknown"
	}
	value, _ = kappnavConfig.Data["kappnav-sa-name"]
	if len(value) == 0 {
		kappnavConfig.Data["kappnav-sa-name"] = instance.GetName() + "-" + ServiceAccountNameSuffix
	}
	if instance.Spec.Console.EnableOkdFeaturedApp {
		kappnavConfig.Data["okd-console-featured-app"] = "enabled"
	} else {
		kappnavConfig.Data["okd-console-featured-app"] = "disabled"
	}
	if instance.Spec.Console.EnableOkdLauncher {
		kappnavConfig.Data["okd-console-app-launcher"] = "enabled"
	} else {
		kappnavConfig.Data["okd-console-app-launcher"] = "disabled"
	}
}

// CreateUIDeploymentContainers ...
func CreateUIDeploymentContainers(existingContainers []corev1.Container, instance *kappnavv1.Kappnav) []corev1.Container {
	// Extract environment variables from existing containers.
	var apiEnv []corev1.EnvVar = nil
	var uiEnv []corev1.EnvVar = nil
	if existingContainers != nil {
		for _, c := range existingContainers {
			switch containerName := c.Name; containerName {
			case APIContainerName:
				apiEnv = c.Env
			case UIContainerName:
				uiEnv = c.Env
			}
		}
	}
	return []corev1.Container{
		*createContainer(APIContainerName, instance, instance.Spec.AppNavAPI, apiEnv,
			createAPIReadinessProbe(), createAPILivenessProbe(), nil),
		*createContainer(UIContainerName, instance, instance.Spec.AppNavUI, uiEnv,
			createUIReadinessProbe(instance), createUILiveinessProbe(instance), createUIVolumeMount(instance)),
	}
}

// CreateControllerDeploymentContainers ...
func CreateControllerDeploymentContainers(existingContainers []corev1.Container, instance *kappnavv1.Kappnav) []corev1.Container {
	// Extract environment variables from existing containers.
	var apiEnv []corev1.EnvVar = nil
	var controllerEnv []corev1.EnvVar = nil
	if existingContainers != nil {
		for _, c := range existingContainers {
			switch containerName := c.Name; containerName {
			case APIContainerName:
				apiEnv = c.Env
			case ControllerContainerName:
				controllerEnv = c.Env
			}
		}
	}
	return []corev1.Container{
		*createContainer(APIContainerName, instance, instance.Spec.AppNavAPI, apiEnv,
			createAPIReadinessProbe(), createAPILivenessProbe(), nil),
		*createContainer(ControllerContainerName, instance, instance.Spec.AppNavController, controllerEnv,
			createControllerReadinessProbe(), createControllerLivenessProbe(), nil),
	}
}

// CreateUIVolumes ...
func CreateUIVolumes(instance *kappnavv1.Kappnav) []corev1.Volume {
	name := instance.Name + "-" + UIVolumeName
	return []corev1.Volume{
		{
			Name: name,
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: name,
				},
			},
		},
	}
}

// SetKappnavDefaults sets default values on the CR instance
func SetKappnavDefaults(instance *kappnavv1.Kappnav, extension KappnavExtension) {
	if extension != nil {
		extension.ApplyAdditionalDefaults(instance)
	}
	setAPIContainerDefaults(instance)
	setUIContainerDefaults(instance)
	setControllerContainerDefaults(instance)
	setImageDefaults(instance)
	setEnvironmentDefaults(instance)
	setArchitectureDefaults(instance)
	setConsoleDefaults(instance)
}

func setAPIContainerDefaults(instance *kappnavv1.Kappnav) {
	apiConfig := instance.Spec.AppNavAPI
	if apiConfig == nil {
		apiConfig = &kappnavv1.KappnavContainerConfiguration{}
		instance.Spec.AppNavAPI = apiConfig
	}
	setContainerDefaults(apiConfig, DefaultAPIRepository)
}

func setUIContainerDefaults(instance *kappnavv1.Kappnav) {
	uiConfig := instance.Spec.AppNavUI
	if uiConfig == nil {
		uiConfig = &kappnavv1.KappnavContainerConfiguration{}
		instance.Spec.AppNavUI = uiConfig
	}
	setContainerDefaults(uiConfig, DefaultUIRepository)
}

func setControllerContainerDefaults(instance *kappnavv1.Kappnav) {
	controllerConfig := instance.Spec.AppNavController
	if controllerConfig == nil {
		controllerConfig = &kappnavv1.KappnavContainerConfiguration{}
		instance.Spec.AppNavController = controllerConfig
	}
	setContainerDefaults(controllerConfig, DefaultControllerRepository)
}

func setContainerDefaults(containerConfig *kappnavv1.KappnavContainerConfiguration, defaultRepoName kappnavv1.Repository) {
	if len(containerConfig.Repository) == 0 {
		containerConfig.Repository = defaultRepoName
	}
	if len(containerConfig.Tag) == 0 {
		containerConfig.Tag = DefaultTag
	}
	if containerConfig.Resources == nil {
		containerConfig.Resources = &kappnavv1.KappnavResourceConstraints{
			Enabled: false,
		}
	} else {
		if containerConfig.Resources.Enabled {
			if containerConfig.Resources.Requests == nil {
				containerConfig.Resources.Requests = &kappnavv1.Resources{}
			}
			setResourceDefaults(containerConfig.Resources.Requests)
			if containerConfig.Resources.Limits == nil {
				containerConfig.Resources.Limits = &kappnavv1.Resources{}
			}
			setResourceDefaults(containerConfig.Resources.Limits)
		}
	}
}

func setResourceDefaults(resources *kappnavv1.Resources) {
	if len(resources.CPU) == 0 {
		resources.CPU = CPUConstraintDefault
	}
	if len(resources.Memory) == 0 {
		resources.Memory = MemoryConstraintDefault
	}
}

func setImageDefaults(instance *kappnavv1.Kappnav) {
	image := instance.Spec.Image
	if image == nil {
		image = &kappnavv1.KappnavImageConfiguration{}
		instance.Spec.Image = image
	}
	if len(image.PullPolicy) == 0 {
		image.PullPolicy = corev1.PullAlways
	}
	if image.PullSecrets == nil {
		image.PullSecrets = []string{}
	}
}

func setEnvironmentDefaults(instance *kappnavv1.Kappnav) {
	env := instance.Spec.Env
	if env == nil {
		env = &kappnavv1.Environment{}
		instance.Spec.Env = env
	}
	if len(env.KubeEnv) == 0 {
		env.KubeEnv = KubeEnvDefault
	}
}

func setArchitectureDefaults(instance *kappnavv1.Kappnav) {
	arch := instance.Spec.Arch
	if arch == nil {
		arch = &kappnavv1.Architecture{}
		instance.Spec.Arch = arch
	}
	if len(arch.Amd64) == 0 {
		arch.Amd64 = kappnavv1.NoPreference
	}
	if len(arch.Ppc64le) == 0 {
		arch.Ppc64le = kappnavv1.NoPreference
	}
	if len(arch.S390x) == 0 {
		arch.S390x = kappnavv1.NoPreference
	}
}

func setConsoleDefaults(instance *kappnavv1.Kappnav) {
	console := instance.Spec.Console
	if console == nil {
		console = &kappnavv1.KappnavConsoleConfiguration{
			EnableOkdFeaturedApp: true,
			EnableOkdLauncher:    true,
		}
		instance.Spec.Console = console
	}
}

func createContainer(name string, instance *kappnavv1.Kappnav,
	containerConfig *kappnavv1.KappnavContainerConfiguration,
	existingEnv []corev1.EnvVar,
	readinessProbe *corev1.Probe,
	livenessProbe *corev1.Probe,
	volumeMount *corev1.VolumeMount) *corev1.Container {
	container := &corev1.Container{
		Name:            name,
		Image:           string(containerConfig.Repository) + ":" + string(containerConfig.Tag),
		ImagePullPolicy: instance.Spec.Image.PullPolicy,
		Env: []corev1.EnvVar{
			{
				Name:  "KAPPNAV_CR_NAME",
				Value: instance.Name,
			},
			{
				Name:  "KAPPNAV_CONFIG_NAMESPACE",
				Value: instance.Namespace,
			},
			{
				Name:  "KUBE_ENV",
				Value: string(instance.Spec.Env.KubeEnv),
			},
		},
		ReadinessProbe: readinessProbe,
		LivenessProbe:  livenessProbe,
	}
	// Copy custom environment variable settings.
	if existingEnv != nil {
		for _, envVar := range existingEnv {
			if envVar.Name != "KAPPNAV_CR_NAME" &&
				envVar.Name != "KAPPNAV_CONFIG_NAMESPACE" &&
				envVar.Name != "KUBE_ENV" {
				container.Env = append(container.Env, envVar)
			}
		}
	}
	// Add volume mount if specified.
	if volumeMount != nil {
		container.VolumeMounts = []corev1.VolumeMount{*volumeMount}
	}
	// Apply resource constraints if enabled.
	if containerConfig.Resources.Enabled {
		container.Resources = corev1.ResourceRequirements{
			Limits:   corev1.ResourceList{},
			Requests: corev1.ResourceList{},
		}
		limits := containerConfig.Resources.Limits
		cpuLimit, err := resource.ParseQuantity(limits.CPU)
		if err == nil {
			container.Resources.Limits[corev1.ResourceCPU] = cpuLimit
		}
		memoryLimit, err := resource.ParseQuantity(limits.Memory)
		if err == nil {
			container.Resources.Limits[corev1.ResourceMemory] = memoryLimit
		}
		requests := containerConfig.Resources.Requests
		cpuRequest, err := resource.ParseQuantity(requests.CPU)
		if err == nil {
			container.Resources.Requests[corev1.ResourceCPU] = cpuRequest
		}
		memoryRequest, err := resource.ParseQuantity(requests.Memory)
		if err == nil {
			container.Resources.Requests[corev1.ResourceMemory] = memoryRequest
		}
	}
	setContainerSecurity(container)
	return container
}

func createAPIReadinessProbe() *corev1.Probe {
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/kappnav/health",
				Scheme: corev1.URISchemeHTTPS,
				Port:   intstr.FromInt(9443),
			},
		},
		InitialDelaySeconds: 60,
		PeriodSeconds:       15,
		FailureThreshold:    6,
	}
}

func createAPILivenessProbe() *corev1.Probe {
	probe := createAPIReadinessProbe()
	probe.InitialDelaySeconds = 120
	return probe
}

func createUIReadinessProbe(instance *kappnavv1.Kappnav) *corev1.Probe {
	kubeEnv := instance.Spec.Env.KubeEnv
	var scheme corev1.URIScheme
	if IsMinikubeEnv(kubeEnv) {
		scheme = corev1.URISchemeHTTP
	} else {
		scheme = corev1.URISchemeHTTPS
	}
	return &corev1.Probe{
		Handler: corev1.Handler{
			HTTPGet: &corev1.HTTPGetAction{
				Path:   "/health",
				Scheme: scheme,
				Port:   intstr.FromInt(3000),
			},
		},
		InitialDelaySeconds: 20,
		PeriodSeconds:       10,
		FailureThreshold:    6,
	}
}

func createUILiveinessProbe(instance *kappnavv1.Kappnav) *corev1.Probe {
	probe := createUIReadinessProbe(instance)
	probe.InitialDelaySeconds = 40
	probe.PeriodSeconds = 30
	return probe
}

func createUIVolumeMount(instance *kappnavv1.Kappnav) *corev1.VolumeMount {
	volumeMount := &corev1.VolumeMount{
		MountPath: UIVolumeMountPath,
		Name:      instance.Name + "-" + UIVolumeName,
	}
	return volumeMount
}

func createControllerReadinessProbe() *corev1.Probe {
	return &corev1.Probe{
		Handler: corev1.Handler{
			Exec: &corev1.ExecAction{
				Command: []string{
					"/bin/bash",
					"-c",
					"testcntlr.sh",
				},
			},
		},
		InitialDelaySeconds: 30,
		PeriodSeconds:       5,
		FailureThreshold:    6,
	}
}

func createControllerLivenessProbe() *corev1.Probe {
	probe := createControllerReadinessProbe()
	probe.InitialDelaySeconds = 120
	probe.PeriodSeconds = 30
	return probe
}

func setContainerSecurity(container *corev1.Container) {
	f := false
	container.SecurityContext = &corev1.SecurityContext{
		Privileged:               &f,
		ReadOnlyRootFilesystem:   &f,
		AllowPrivilegeEscalation: &f,
		Capabilities: &corev1.Capabilities{
			Drop: []corev1.Capability{"ALL"},
		},
	}
}

func setPodSecurity(pts *corev1.PodTemplateSpec) {
	pts.Spec.HostNetwork = false
	pts.Spec.HostPID = false
	pts.Spec.HostIPC = false
	t := true
	user := int64(1001)
	pts.Spec.SecurityContext = &corev1.PodSecurityContext{
		RunAsNonRoot: &t,
		RunAsUser:    &user,
	}
}

// IsMinikubeEnv ...
func IsMinikubeEnv(kubeEnv string) bool {
	return kubeEnv == "minikube" || kubeEnv == "k8s"
}

//
// Functions for accessing and updating status on the CR
//

// GetCondition ...
func GetCondition(conditionType kappnavv1.StatusConditionType, status *kappnavv1.KappnavStatus) *kappnavv1.StatusCondition {
	for i := range status.Conditions {
		if status.Conditions[i].Type == conditionType {
			return &status.Conditions[i]
		}
	}
	return nil
}

// SetCondition ...
func SetCondition(condition kappnavv1.StatusCondition, status *kappnavv1.KappnavStatus) {
	for i := range status.Conditions {
		if status.Conditions[i].Type == condition.Type {
			status.Conditions[i] = condition
			return
		}
	}
	status.Conditions = append(status.Conditions, condition)
}