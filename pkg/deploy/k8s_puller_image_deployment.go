package deploy

import (
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/cmd/operator_osd/config"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KubernetesPullerImageDeployment() *appsv1.Deployment{
	return &appsv1.Deployment{
		TypeMeta:   metav1.TypeMeta{
			Kind: "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "kubernetes-image-puller",
			},
			Name: "kubernetes-image-puller",
		},
		Spec: appsv1.DeploymentSpec {
			Replicas: int32Ptr(1),
			RevisionHistoryLimit: int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "kubernetes-image-puller",
				},
				MatchExpressions: nil,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app" : "kubernetes-image-puller",
					},
				},
				Spec: v1.PodSpec{
					ServiceAccountName: "k8s-image-puller",
					Containers: []v1.Container{
						{
							Name: "kubernetes-image-puller",
							Image: config.TestHarnessConfig.KubernetesImagePuller.Image,
							ImagePullPolicy: "IfNotPresent",
							EnvFrom: []v1.EnvFromSource{
								{
									ConfigMapRef: &v1.ConfigMapEnvSource{
										LocalObjectReference: v1.LocalObjectReference{
											Name: DefaultConfigMapName,
										},
									},
								},
							},
						},
					},
				},
			},
			Strategy: appsv1.DeploymentStrategy{
				Type: "Recreate",
			},
		},
	}
}
