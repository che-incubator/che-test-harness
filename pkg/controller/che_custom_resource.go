package controller

import (
	"errors"
	"github.com/che-incubator/che-test-harness/cmd/che/config"
	"time"

	"github.com/che-incubator/che-test-harness/pkg/deploy"
	"github.com/che-incubator/che-test-harness/pkg/monitors"
	"github.com/che-incubator/che-test-harness/pkg/monitors/metadata"
	orgv1 "github.com/eclipse/che-operator/pkg/apis/org/v1"
	"go.uber.org/zap"
)

const (
	CheKind = "checlusters"
)

// GetCustomResource make an api request to K8s API to get Che Cluster
func (c *TestHarnessController) GetCustomResource() (*orgv1.CheCluster, error) {
	result := orgv1.CheCluster{}

	err := c.kubeClient.KubeRest().
		Get().
		Namespace(metadata.Namespace.Name).
		Resource(CheKind).
		Name(config.TestHarnessConfig.Flavor).
		Do().
		Into(&result)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateCustomResource make an api request to K8s API to delete Che Cluster
func (c *TestHarnessController) CreateCustomResource() (err error) {
	result := orgv1.CheCluster{}
	cheCluster := deploy.CreateCodeReadyCluster()

	err = c.kubeClient.KubeRest().
		Post().
		Namespace(metadata.Namespace.Name).
		Resource(CheKind).
		Name(config.TestHarnessConfig.Flavor).
		Body(cheCluster).
		Do().
		Into(&result)

	if err != nil {
		c.Logger.Error("Error on create custom resource", zap.Error(err))
	}

	return err
}

// WatchCustomResource wait to deploy all che/crw pods
func (c *TestHarnessController) WatchCustomResource(status string) (deployed bool, err error) {
	var clusterStarted = time.Now()
	timeout := time.After(15 * time.Minute)
	tick := time.Tick(1 * time.Second)

	stopCh := make(chan struct{})
	defer close(stopCh)
	monitor := monitors.NewPodStartupDataMonitor(c.kubeClient.Kube())
	go func() {
		if err := monitor.DescribeEvents(stopCh); err != nil {
			panic(err)
		}
	}()

	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			customResource, _ := c.GetCustomResource()
			if customResource.Status.CheClusterRunning != status {
			} else {
				metadata.Instance.ClusterTimeUp = time.Since(clusterStarted).Seconds()
				return true, nil
			}
		}
	}
}

// DeleteCustomResource make an api request to K8s API to delete Che Cluster
func (c TestHarnessController) DeleteCustomResource() (err error) {
	err = c.kubeClient.KubeRest().
		Delete().
		Namespace(metadata.Namespace.Name).
		Resource(CheKind).
		Name(config.TestHarnessConfig.Flavor).
		Do().
		Error()

	if err != nil {
		c.Logger.Error("Error on delete custom resource", zap.Error(err))

		return err
	}

	return nil
}
