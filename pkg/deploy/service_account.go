package deploy

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func PullerImageServiceAccount() *v1.ServiceAccount {
	return &v1.ServiceAccount{
		TypeMeta : metav1.TypeMeta{
			Kind: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "k8s-image-puller",
		},
	}
}
