package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KappnavSpec defines the desired state of Kappnav
// +k8s:openapi-gen=true
type KappnavSpec struct {
	AppNavAPI        *KappnavContainerConfiguration        `json:"appNavAPI,omitempty"`
	AppNavController *KappnavContainerConfiguration        `json:"appNavController,omitempty"`
	AppNavUI         *KappnavServiceContainerConfiguration `json:"appNavUI,omitempty"`
	AppNavInit       *KappnavContainerConfiguration        `json:"appNavInit,omitempty"`
}

// KappnavContainerConfiguration defines the configuration for a Kappnav container
type KappnavContainerConfiguration struct {
	Repository Repository                 `json:"repository,omitempty"`
	Tag        Tag                        `json:"tag,omitempty"`
	Resources  KappnavResourceConstraints `json:"resources,omitempty"`
}

// KappnavResourceConstraints defines resource constraints for a Kappnav container
type KappnavResourceConstraints struct {
	Enabled  bool      `json:"enabled,omitempty"`
	Requests Resources `json:"requests,omitempty"`
	Limits   Resources `json:"limits,omitempty"`
}

// Resources ...
type Resources struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// Repository ...
type Repository string

// Tag ...
type Tag string

// KappnavServiceContainerConfiguration defines the configuration for a Kappnav container with a service
type KappnavServiceContainerConfiguration struct {
    KappnavContainerConfiguration `json:",inline"`
	Service Service               `json:"service,omitempty"`
}

// Service ...
type Service struct {
	Type ServiceType `json:"type,omitempty"`
}

// ServiceType ...
type ServiceType string

// KappnavStatus defines the observed state of Kappnav
// +k8s:openapi-gen=true
type KappnavStatus struct {
	LastTransitionTime *metav1.Time           `json:"lastTransitionTime,omitempty"`
	LastUpdateTime     *metav1.Time           `json:"lastUpdateTime,omitempty"`
	Reason             string                 `json:"reason,omitempty"`
	Message            string                 `json:"message,omitempty"`
	Status             corev1.ConditionStatus `json:"status,omitempty"`
	Type               StatusConditionType    `json:"type,omitempty"`
}

// StatusConditionType ...
type StatusConditionType string

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Kappnav is the Schema for the kappnavs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
type Kappnav struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KappnavSpec   `json:"spec,omitempty"`
	Status KappnavStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// KappnavList contains a list of Kappnav
type KappnavList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Kappnav `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Kappnav{}, &KappnavList{})
}
