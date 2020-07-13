package controller

import (
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/cmd/operator_osd/config"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/deploy"
)

const (
	KubernetesPullerImageLabel = "test=daemonset-test"
)
// DeployKubernetesPullerImage Creates all configs and deploy kubernetes image puller in cluster
func (c *TestHarnessController) DeployKubernetesPullerImage() error {
	c.Logger.Info("Starting to deploy Kubernetes Puller deploy...")

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

	_, err := c.WatchPodStartup(config.TestHarnessConfig.KubernetesImagePuller.Namespace, KubernetesPullerImageLabel, "")

	if err != nil {
		return err
	}

	c.Logger.Info("Kubernetes Puller Image was deployed successfully...")

	return nil
}

// CreateKubernetesPullerImageServiceAccount create a service account for kubernetes image puller
func (c *TestHarnessController) CreateKubernetesPullerImageServiceAccount() error {
	sa := deploy.PullerImageServiceAccount()
	_, err := c.kubeClient.Kube().CoreV1().ServiceAccounts(config.TestHarnessConfig.KubernetesImagePuller.Namespace).Create(sa)

	return err
}

// CreateKubernetesPullerImageRole create roles for kubernetes image puller
func (c *TestHarnessController) CreateKubernetesPullerImageRole() error {
	role := deploy.KubernetesPullerImageRole()
	_, err := c.kubeClient.Kube().RbacV1().Roles(config.TestHarnessConfig.KubernetesImagePuller.Namespace).Create(role)

	return err
}
// CreateKubernetesPullerImageRoleBinding create roles binding for kubernetes image puller
func (c *TestHarnessController) CreateKubernetesPullerImageRoleBinding() error {
	roleBinding := deploy.KubernetesPullerImageRoleBinding()
	_, err := c.kubeClient.Kube().RbacV1().RoleBindings(config.TestHarnessConfig.KubernetesImagePuller.Namespace).Create(roleBinding)

	return err
}

// CreateKubernetesPullerImageConfigMap create config maps for kubernetes image puller
func (c *TestHarnessController) CreateKubernetesPullerImageConfigMap() error {
	cfg := deploy.KubernetesPullerImageConfigMap()
	_, err := c.kubeClient.Kube().CoreV1().ConfigMaps(config.TestHarnessConfig.KubernetesImagePuller.Namespace).Create(cfg)

	return err
}

// CreateKubernetesPullerImageDeployment create deployment for kubernetes image puller
func (c *TestHarnessController) CreateKubernetesPullerImageDeployment() error {
	deployment := deploy.KubernetesPullerImageDeployment()
	_, err := c.kubeClient.Kube().AppsV1().Deployments(config.TestHarnessConfig.KubernetesImagePuller.Namespace).Create(deployment)

	return err
}
