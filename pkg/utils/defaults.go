package utils

import (
	kappnavv1 "github.com/kappnav/operator/pkg/apis/kappnav/v1"
	"io/ioutil"
	"sigs.k8s.io/yaml"
)

// SetKappnavDefaults sets default values on the CR instance
func SetKappnavDefaults(instance *kappnavv1.Kappnav, extension KappnavExtension) error {
	defaults, err := getDefaults()
	if err != nil {
		return err
	}
	if extension != nil {
		extension.ApplyAdditionalDefaults(instance, defaults)
	}
	setAPIContainerDefaults(instance, defaults)
	setUIContainerDefaults(instance, defaults)
	setControllerContainerDefaults(instance, defaults)
	setImageDefaults(instance, defaults)
	setEnvironmentDefaults(instance, defaults)
	return nil
}

func getDefaults() (*kappnavv1.Kappnav, error) {
	// Read default values file
	fData, err := ioutil.ReadFile("deploy/default_values.yaml")
	if err != nil {
		return nil, err
	}
	defaults := &kappnavv1.Kappnav{}
	err = yaml.Unmarshal(fData, defaults)
	if err != nil {
		return nil, err
	}
	return defaults, nil
}

func setAPIContainerDefaults(instance *kappnavv1.Kappnav, defaults *kappnavv1.Kappnav) {
	apiConfig := instance.Spec.AppNavAPI
	if apiConfig == nil {
		instance.Spec.AppNavAPI = defaults.Spec.AppNavAPI
	} else {
		setContainerDefaults(apiConfig, defaults.Spec.AppNavAPI)
	}
}

func setUIContainerDefaults(instance *kappnavv1.Kappnav, defaults *kappnavv1.Kappnav) {
	uiConfig := instance.Spec.AppNavUI
	if uiConfig == nil {
		instance.Spec.AppNavUI = defaults.Spec.AppNavUI
	} else {
		setContainerDefaults(uiConfig, defaults.Spec.AppNavUI)
	}
}

func setControllerContainerDefaults(instance *kappnavv1.Kappnav, defaults *kappnavv1.Kappnav) {
	controllerConfig := instance.Spec.AppNavController
	if controllerConfig == nil {
		instance.Spec.AppNavController = defaults.Spec.AppNavController
	} else {
		setContainerDefaults(controllerConfig, defaults.Spec.AppNavController)
	}
}

func setContainerDefaults(containerConfig *kappnavv1.KappnavContainerConfiguration,
	defaultContainerConfig *kappnavv1.KappnavContainerConfiguration) {
	if len(containerConfig.Repository) == 0 {
		containerConfig.Repository = defaultContainerConfig.Repository
	}
	if len(containerConfig.Tag) == 0 {
		containerConfig.Tag = defaultContainerConfig.Tag
	}
	if containerConfig.Resources == nil {
		containerConfig.Resources = defaultContainerConfig.Resources
	} else {
		if containerConfig.Resources.Enabled {
			if containerConfig.Resources.Requests == nil {
				containerConfig.Resources.Requests = defaultContainerConfig.Resources.Requests
			} else {
				setResourceDefaults(containerConfig.Resources.Requests, defaultContainerConfig.Resources.Requests)
			}
			if containerConfig.Resources.Limits == nil {
				containerConfig.Resources.Limits = defaultContainerConfig.Resources.Limits
			} else {
				setResourceDefaults(containerConfig.Resources.Limits, defaultContainerConfig.Resources.Limits)
			}
		}
	}
}

func setResourceDefaults(resources *kappnavv1.Resources, defaultResources *kappnavv1.Resources) {
	if len(resources.CPU) == 0 {
		resources.CPU = defaultResources.CPU
	}
	if len(resources.Memory) == 0 {
		resources.Memory = defaultResources.Memory
	}
}

func setImageDefaults(instance *kappnavv1.Kappnav, defaults *kappnavv1.Kappnav) {
	image := instance.Spec.Image
	if image == nil {
		instance.Spec.Image = defaults.Spec.Image
	} else {
		if len(image.PullPolicy) == 0 {
			image.PullPolicy = defaults.Spec.Image.PullPolicy
		}
		if image.PullSecrets == nil {
			image.PullSecrets = defaults.Spec.Image.PullSecrets
		}
	}
}

func setEnvironmentDefaults(instance *kappnavv1.Kappnav, defaults *kappnavv1.Kappnav) {
	env := instance.Spec.Env
	if env == nil {
		instance.Spec.Env = defaults.Spec.Env
	} else {
		if len(env.KubeEnv) == 0 {
			env.KubeEnv = defaults.Spec.Env.KubeEnv
		}
	}
}