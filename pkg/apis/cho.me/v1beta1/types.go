package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FooServer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FooServerSpec   `json:"spec,omitempty"`
	Status FooServerStatus `json:"status,omitempty"`
}

type FooServerSpec struct {
	DeploymentSpec FooServerDeploymentSpec `json:"deploymentSpec"`
	ServiceSpec    FooServerServiceSpec    `json:"serviceSpec"`
	SecretSpec     FooServerSecretSpec     `json:"secretSpec,omitempty"`
}

type FooServerDeploymentSpec struct {
	Name             string `json:"name"`
	PodReplicas      *int32 `json:"podReplicas"`
	PodContainerPort *int32 `json:"podContainerPort"`
}

type FooServerServiceSpec struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Port       *int32 `json:"port"`
	TargetPort *int32 `json:"targetPort"`
	NodePort   *int32 `json:"nodePort"`
}

type FooServerSecretSpec struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type FooServerStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FooServerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []FooServer `json:"items"`
}
