package controller

import (
	"sync"

	"github.com/eclipse/che-operator/pkg/controller/che"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/client"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/controller/logger"
	"go.uber.org/zap"
)

// TestHarnessController useful to add all kubernetes objects to cluster.
type TestHarnessController struct {
	sync.Mutex
	kubeClient *client.K8sClient
	Logger     logger.Log
}

// NewTestHarnessController creates a new TestHarnessController from a given client.
func NewTestHarnessController(c *client.K8sClient) *TestHarnessController {
	return &TestHarnessController{
		kubeClient: c,
		Logger:     logger.Zap,
	}
}

func (c *TestHarnessController) RunTestHarness() bool {
	if err := c.DeployKubernetesPullerImage(); err != nil {
		c.Logger.Panic("Failed to deploy Kubernetes Puller Image in cluster", zap.Error(err))
	}

	c.Logger.Info("Generating Custom Resource in cluster")
	//Create a new Code Ready Workspaces Custom resources into a giving namespace.
	if err := c.CreateCustomResource(); err != nil {
		c.Logger.Panic("Failed to create custom resources in cluster", zap.Error(err))
	}

	c.Logger.Info("Successfully created CodeReady Custom Resources")

	// Check If all kubernetes objects for code ready workspaces are created in cluster
	// !Timeout is 15 minutes
	c.Logger.Info("Starting to check Code Ready Cluster if is available")
	deploy, _ := c.WatchCustomResource(che.AvailableStatus)

	return deploy
}
