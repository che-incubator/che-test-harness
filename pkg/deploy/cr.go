package deploy

import (
	"errors"
	"github.com/che-incubator/che-test-harness/pkg/deploy/che"
	"github.com/che-incubator/che-test-harness/pkg/deploy/context"
	"github.com/che-incubator/che-test-harness/pkg/monitors"
	"time"

	orgv1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	"go.uber.org/zap"
)

// GetCustomResource make an models request to K8s API to get Che Cluster
func (c *TestHarnessController) GetCustomResource() (*orgv1.CheCluster, error) {
	result := orgv1.CheCluster{}

	err := c.kubeClient.KubeRest().
		Get().
		Namespace(context.TestInstance.Setup.CheNamespace).
		Resource(context.CustomResources).
		Name(context.CheCustomResourceName).
		Do().
		Into(&result)
	if err != nil {

		return nil, err
	}

	return &result, nil
}

// CreateCustomResource make an models request to K8s API to delete Che Cluster
func (c *TestHarnessController) CreateCustomResource() (err error) {
	result := orgv1.CheCluster{}
	cheCluster := che.CreateEclipseCheCluster()

	err = c.kubeClient.KubeRest().
		Post().
		Namespace(context.TestInstance.Setup.CheNamespace).
		Resource(context.CustomResources).
		Name(context.CheCustomResourceName).
		Body(cheCluster).
		Do().
		Into(&result)

	if err != nil {
		c.Logger.Error("Error on create custom resource", zap.Error(err))
	}

	return err
}

// WatchCustomResource wait to deploy all performance/crw pods
func (c *TestHarnessController) WaitForCheToBeReady(status string) (deployed bool, err error) {
	timeout := time.After(10 * time.Minute)
	tick := time.Tick(1 * time.Second)

	stopCh := make(chan struct{})
	defer close(stopCh)
	monitor, _ := monitors.NewMonitor(c.kubeClient.Kube())
	go func() {
		if err := monitor.DescribeEvents(stopCh, ""); err != nil {
			panic(err)
		}
	}()

	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			customResource, _ := c.GetCustomResource()
			if customResource.Status.CheClusterRunning == status {
				context.TestInstance.Metrics.CheClusterUpTime = time.Since(customResource.CreationTimestamp.Time).Seconds()
				return true, nil
			}
		}
	}
}

// DeleteCustomResource make an models request to K8s API to delete Che Cluster
func (c *TestHarnessController) DeleteCustomResource() (err error) {
	err = c.kubeClient.KubeRest().
		Delete().
		Namespace(context.TestInstance.Setup.CheNamespace).
		Resource(context.CustomResources).
		Name(context.CheCustomResourceName).
		Do().
		Error()

	if err != nil {
		c.Logger.Error("Error on delete custom resource", zap.Error(err))

		return err
	}

	return nil
}

