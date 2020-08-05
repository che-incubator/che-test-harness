package controller

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Create a new kubernetes image puller
// !TODO Make creation of namespace configurable. Don't create namespace only for k8s-image-puller. Same for Namespace deletion
func (c *TestHarnessController) CreateNamespace() error {
	_, err := c.kubeClient.Kube().CoreV1().Namespaces().Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: KubernetesImgPullerNS,
		},
	})
	return err
}

// Delete a kubernetes image puller namespace
func (c *TestHarnessController) DeleteNamespace() error {
	opts := metav1.DeleteOptions{}

	err := c.kubeClient.Kube().CoreV1().Namespaces().Delete(KubernetesImgPullerNS, &opts)
	return err
}
