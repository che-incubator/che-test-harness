package controller

import (
	"github.com/che-incubator/che-test-harness/pkg/deploy"
)

const (
	KubernetesPullerImageLabel = "test=daemonset-test"
	KubernetesImgPullerNS      = "k8s-image-puller"
)

// DeployKubernetesPullerImage Creates all configs and deploy kubernetes image puller in cluster
func (c *TestHarnessController) DeployKubernetesPullerImage() error {
	c.Logger.Info("Starting to deploy Kubernetes Puller")

	if err := c.CreateNamespace(); err != nil {
		return err
	}
	c.Logger.Info("Successfully created namespace for Kubernetes Puller Image...")

	if err := c.CreateKubernetesPullerImageServiceAccount(); err != nil {
		return err
	}
	c.Logger.Info("Successfully created service account for Kubernetes Puller Image...")

	if err := c.CreateKubernetesPullerImageRole(); err != nil {
		return err
	}
	c.Logger.Info("Successfully created service roles for Kubernetes Puller Image...")

	if err := c.CreateKubernetesPullerImageRoleBinding(); err != nil {
		return err
	}
	c.Logger.Info("Successfully created roleBinding for Kubernetes Puller Image...")

	if err := c.CreateKubernetesPullerImageConfigMap(); err != nil {
		return err
	}
	c.Logger.Info("Successfully created configMaps for Kubernetes Puller Image...")

	if err := c.CreateKubernetesPullerImageDeployment(); err != nil {
		return err
	}

	_, err := c.WatchPodStartup(KubernetesImgPullerNS, KubernetesPullerImageLabel, "")

	if err != nil {
		return err
	}

	c.Logger.Info("Kubernetes Puller Image was deployed successfully...")

	return nil
}

// CreateKubernetesPullerImageServiceAccount create a service account for kubernetes image puller
func (c *TestHarnessController) CreateKubernetesPullerImageServiceAccount() error {
	sa := deploy.PullerImageServiceAccount()
	_, err := c.kubeClient.Kube().CoreV1().ServiceAccounts(KubernetesImgPullerNS).Create(sa)

	return err
}

// CreateKubernetesPullerImageRole create roles for kubernetes image puller
func (c *TestHarnessController) CreateKubernetesPullerImageRole() error {
	role := deploy.KubernetesPullerImageRole()
	_, err := c.kubeClient.Kube().RbacV1().Roles(KubernetesImgPullerNS).Create(role)

	return err
}

// CreateKubernetesPullerImageRoleBinding create roles binding for kubernetes image puller
func (c *TestHarnessController) CreateKubernetesPullerImageRoleBinding() error {
	roleBinding := deploy.KubernetesPullerImageRoleBinding()
	_, err := c.kubeClient.Kube().RbacV1().RoleBindings(KubernetesImgPullerNS).Create(roleBinding)

	return err
}

// CreateKubernetesPullerImageConfigMap create secrets maps for kubernetes image puller
func (c *TestHarnessController) CreateKubernetesPullerImageConfigMap() error {
	cfg := deploy.KubernetesPullerImageConfigMap()
	_, err := c.kubeClient.Kube().CoreV1().ConfigMaps(KubernetesImgPullerNS).Create(cfg)

	return err
}

// CreateKubernetesPullerImageDeployment create deployment for kubernetes image puller
func (c *TestHarnessController) CreateKubernetesPullerImageDeployment() error {
	deployment := deploy.KubernetesPullerImageDeployment()
	_, err := c.kubeClient.Kube().AppsV1().Deployments(KubernetesImgPullerNS).Create(deployment)

	return err
}
