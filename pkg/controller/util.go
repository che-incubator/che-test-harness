package controller

import (
	"errors"
	"fmt"
	"time"

	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/monitors/metadata"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
)

func (c *TestHarnessController) WatchPodStartup(namespace string, label string, workspaceStack string) (deployed bool, err error) {
	timeout := time.After(15 * time.Minute)
	tick := time.Tick(1 * time.Second)

	for {
		select {
		case <-timeout:
			return false, errors.New("timed out")
		case <-tick:
			desc := c.WaitForPodBySelectorRunning(namespace, label, 180, workspaceStack)
			if desc != nil {
			} else {
				return true, nil
			}
		}
	}
}

// return a condition function that indicates whether the given pod is
// currently running
func (c *TestHarnessController) isPodRunning(podName, namespace string, workspaceStack string) wait.ConditionFunc {
	return func() (bool, error) {
		pod, err := c.kubeClient.Kube().CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{IncludeUninitialized: true})
		isReady := podutil.IsPodAvailable(pod, 0, metav1.Time{})

		if isReady {
			_ = metadata.WriteWorkspaceMeasureTimesToMetadata(pod, workspaceStack)

			return true, nil
		}

		if err != nil {
			return false, err
		}
		return false, nil
	}
}

// Poll up to timeout seconds for pod to enter running state.
// Returns an error if the pod never enters the running state.
func (c *TestHarnessController) waitForPodRunning(namespace, podName string, timeout time.Duration, workspaceStack string) error {
	return wait.PollImmediate(time.Second, timeout, c.isPodRunning(podName, namespace, workspaceStack))
}

// Returns the list of currently scheduled or running pods in `namespace` with the given selector
func (c *TestHarnessController) ListPods(namespace, selector string) (*v1.PodList, error) {
	listOptions := metav1.ListOptions{IncludeUninitialized: true, LabelSelector: selector}
	podList, err := c.kubeClient.Kube().CoreV1().Pods(namespace).List(listOptions)

	if err != nil {
		return nil, err
	}
	return podList, nil
}

// Wait up to timeout seconds for all pods in 'namespace' with given 'selector' to enter running state.
// Returns an error if no pods are found or not all discovered pods enter running state.
func (c *TestHarnessController) WaitForPodBySelectorRunning(namespace, selector string, timeout int, workspaceStack string) error {
	podList, err := c.ListPods(namespace, selector)
	if err != nil {
		return err
	}
	if len(podList.Items) == 0 {
		c.Logger.Warn("Pod not created yet with selector " + selector + " in namespace " + namespace)

		return fmt.Errorf("Pod not created yet in %s with selector %s", namespace, selector)
	}
	for _, pod := range podList.Items {
		c.Logger.Info("Pod " + pod.Name + " created in namespace " + namespace)
		if err := c.waitForPodRunning(namespace, pod.Name, time.Duration(timeout)*time.Second, workspaceStack); err != nil {
			return err
		}
	}

	return nil
}
