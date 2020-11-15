package deploy

import (
	"github.com/che-incubator/che-test-harness/internal/logger"
	"github.com/che-incubator/che-test-harness/pkg/client"
	"github.com/eclipse/che-operator/pkg/controller/che"
	"go.uber.org/zap"
	"sync"
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
	c.Logger.Info("Generating Custom Resource in cluster")
	//Create a new Eclipse Che Custom resources into a giving namespace.
	if err := c.CreateCustomResource(); err != nil {
		c.Logger.Panic("Failed to create custom resources in cluster", zap.Error(err))
	}

	c.Logger.Info("Successfully created Eclipse Che Custom Resources")

	// Check If all kubernetes objects for eclipse performance are created in cluster
	// !Timeout is 10 minutes
	c.Logger.Info("Starting to check if Eclipse Che Cluster is available")
	deploy, _ := c.WaitForCheToBeReady(che.AvailableStatus)

	return deploy
}
