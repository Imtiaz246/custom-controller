package controller

import (
	customv1beta1 "github.com/imtiaz246/custom-controller/pkg/apis/cho.me/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func newDeployment(fooServer *customv1beta1.FooServer) *appsv1.Deployment {
	labels := map[string]string{
		"app":        "api-server-deployment",
		"controller": fooServer.Name,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fooServer.Name,
			Namespace: fooServer.Namespace,
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(fooServer, customv1beta1.SchemeGroupVersion.WithKind("FooServer")),
			},
			Labels: labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: fooServer.Spec.DeploymentSpec.PodReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "book-server-container",
							Image: "imtiazcho/book-server",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3000,
								},
							},
							LivenessProbe: &corev1.Probe{
								ProbeHandler: corev1.ProbeHandler{
									HTTPGet: &corev1.HTTPGetAction{
										Port: intstr.IntOrString{
											IntVal: 3000,
											StrVal: "3000",
										},
									},
								},
								InitialDelaySeconds: 10,
								TimeoutSeconds:      1,
								PeriodSeconds:       10,
								FailureThreshold:    3,
							},
							Env: []corev1.EnvVar{
								{
									Name:  "ADMIN_USERNAME",
									Value: "admin",
									//ValueFrom: &corev1.EnvVarSource{
									//	SecretKeyRef: &corev1.SecretKeySelector{
									//		LocalObjectReference: corev1.LocalObjectReference{
									//			Name: "fooServerSecret",
									//		},
									//		Key: "ADMIN_USERNAME",
									//	},
									//},
								},
								{
									Name:  "ADMIN_PASSWORD",
									Value: "password",
									//ValueFrom: &corev1.EnvVarSource{
									//	SecretKeyRef: &corev1.SecretKeySelector{
									//		LocalObjectReference: corev1.LocalObjectReference{
									//			Name: "fooServerSecret",
									//		},
									//		Key: "ADMIN_PASSWORD",
									//	},
									//},
								},
							},
						},
					},
				},
			},
		},
	}
}
