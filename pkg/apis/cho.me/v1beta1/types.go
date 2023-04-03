package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FooServer struct {
	metav1.TypeMeta   `json:"inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FooServerSpec   `json:"spec"`
	Status FooServerStatus `json:"status"`
}

type FooServerSpec struct {
	DeploymentSpec FooServerDeploymentSpec `json:"deploymentSpec"`
	ServiceSpec    FooServerServiceSpec    `json:"serviceSpec"`
	SecretSpec     FooServerSecretSpec     `json:"secretSpec"`
}

type FooServerDeploymentSpec struct {
	DeploymentName   string `json:"deploymentName"`
	PodReplicas      *int32 `json:"podReplicas"`
	PodContainerPort *int32 `json:"podContainerPort"`
}

type FooServerServiceSpec struct {
	ServiceName       string `json:"serviceName"`
	ServiceType       string `json:"serviceType"`
	ServicePort       *int   `json:"servicePort"`
	ServiceTargetPort *int   `json:"serviceTargetPort"`
}

type FooServerSecretSpec struct {
	SecretName   string `json:"secretName"`
	UsernameData string `json:"usernameData"`
	PasswordData string `json:"passwordData"`
}

type FooServerStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type FooServerList struct {
	metav1.TypeMeta   `json:"inline"`
	metav1.ObjectMeta `json:"metadata"`

	Items []FooServer `json:"items"`
}
