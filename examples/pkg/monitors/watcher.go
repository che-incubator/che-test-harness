package monitors

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	testContext "github.com/che-incubator/che-test-harness/pkg/deploy/context"
	"github.com/golang/glog"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	client "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/pager"
	"k8s.io/kubernetes/pkg/kubelet/events"
)

var (
	purgeAfter            time.Duration
	includeInitContainers string
)

type podStartupLatencyDataMonitor struct {
	kubeClient client.Interface
	inits      map[string]bool
	pods       map[string]*podMilestones

	// Map of pods marked for deletion waiting to be purged. Just an
	// optimization to avoid full scans.
	toDelete map[string]time.Time

	// Time after which deleted entries are purged. We don't want to delete
	// data right after observing deletions, because we may get events
	// corresponding to the given Pod later, which would result in recreating
	// entry, and a serious memory leak.
	purgeAfter time.Duration

	initDone bool
	sync.Mutex

	ready bool

	readyAt time.Time
}

// NewMonitor start a controller to get pod times
func NewMonitor(k8s client.Interface) (*podStartupLatencyDataMonitor, error) {
	inits := make(map[string]bool)
	for _, name := range strings.Split(includeInitContainers, ",") {
		if n := strings.TrimSpace(name); n != "" {
			inits[n] = true
		}
	}

	glog.Infof("starting SLO monitor: initContainers: %v purgeAfter: %s", inits, purgeAfter)

	return &podStartupLatencyDataMonitor{
		kubeClient: k8s,
		inits:      inits,
		pods:       make(map[string]*podMilestones),
		toDelete:   make(map[string]time.Time),
		purgeAfter: purgeAfter,
	}, nil
}

// Descibe all pod events in given namespace for Eclipse Che
func (pm *podStartupLatencyDataMonitor) DescribeEvents(stopCh chan struct{}, workspaceStack string) error {
	pm.initDone = true
	_, controller := cache.NewInformer(&cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			pg := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
				return pm.kubeClient.CoreV1().Pods(testContext.TestInstance.CheNamespace).List(opts)
			}))
			return pg.List(context.Background(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return pm.kubeClient.CoreV1().Pods(testContext.TestInstance.CheNamespace).Watch(options)
		},
	}, new(v1.Pod), 0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				_ = obj.(*v1.Pod)
			},
			DeleteFunc: func(obj interface{}) {},
			UpdateFunc: func(oldObj, newObj interface{}) {
				pod := newObj.(*v1.Pod)
				pm.handlePodUpdate(pod, workspaceStack)
			},
		})
	go controller.Run(stopCh)

	eventSelector := fields.Set{
		"involvedObject.kind": "Pod",
		"source":              "kubelet",
	}.AsSelector().String()

	_, eventcontroller := cache.NewInformer(&cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = eventSelector
			pg := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
				return pm.kubeClient.CoreV1().Events(testContext.TestInstance.CheNamespace).List(opts)
			}))
			return pg.List(context.Background(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = eventSelector
			return pm.kubeClient.CoreV1().Events(testContext.TestInstance.CheNamespace).Watch(options)
		},
	}, new(v1.Event), 0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				e := obj.(*v1.Event)
				pm.handleEvent(e, workspaceStack)
			},
			DeleteFunc: func(obj interface{}) {},
		})
	go eventcontroller.Run(stopCh)

	return nil
}

func (pm *podStartupLatencyDataMonitor) handleEvent(event *v1.Event, workspaceStack string) {
	if !pm.initDone {
		return
	}
	switch event.Reason {
	case events.PullingImage:
		go func() {
			if err := pm.handlePullingImageEvent(event, workspaceStack); err != nil {
				glog.Warningf("failed to process 'PullingImage' event: %v", err)
			}
		}()
	case events.PulledImage:
		go func() {
			if err := pm.handlePulledImageEvent(event, workspaceStack); err != nil {
				glog.Warningf("failed to process 'PulledImage' event: %v", err)
			}
		}()
	}
}

func (pm *podStartupLatencyDataMonitor) handlePullingImageEvent(event *v1.Event, workspaceStack string) error {
	containerName, err := fetchContainerName(event.InvolvedObject)
	if err != nil {
		return fmt.Errorf("failed to fetch container name: %v", err)
	}

	containerMilestones := newContainerMilestones(containerName)
	containerMilestones.pullingAt = event.FirstTimestamp.Time
	containerMilestones.pulling = true

	pm.Lock()
	defer pm.Unlock()

	podKey := getPodKeyFromReference(event.InvolvedObject)
	podMilestones, ok := pm.pods[podKey]
	if !ok {
		podMilestones = newPodMilestonesFromReference(event.InvolvedObject)
		pm.pods[podKey] = podMilestones
	}

	podMilestones.mergeContainer(containerMilestones)
	glog.V(4).Infof("handlePullingImageEvent %q: %s", podKey, podMilestones)

	return pm.updateMetric(podMilestones, workspaceStack)
}

func (pm *podStartupLatencyDataMonitor) handlePulledImageEvent(event *v1.Event, workspaceStack string) error {
	containerName, err := fetchContainerName(event.InvolvedObject)
	if err != nil {
		return fmt.Errorf("failed to fetch container name: %v", err)
	}

	containerMilestones := newContainerMilestones(containerName)
	containerMilestones.alreadyPresent = strings.Contains(event.Message, "already present on machine")
	containerMilestones.pulledAt = event.FirstTimestamp.Time
	containerMilestones.pulled = true

	pm.Lock()
	defer pm.Unlock()

	podKey := getPodKeyFromReference(event.InvolvedObject)
	podMilestones, ok := pm.pods[podKey]
	if !ok {
		podMilestones = newPodMilestonesFromReference(event.InvolvedObject)
		pm.pods[podKey] = podMilestones
	}

	podMilestones.mergeContainer(containerMilestones)
	glog.V(4).Infof("handlePulledImageEvent %q: %s", podKey, podMilestones)

	return pm.updateMetric(podMilestones, workspaceStack)
}

func (pm *podStartupLatencyDataMonitor) handlePodUpdate(pod *v1.Pod, workspaceStack string) {
	if !pm.initDone {
		return
	}

	if pml := newPodMilestonesFromPod(pod); pml.allContainerStarted() {
		go func() {
			ready, readyAt := checkPodAndGetStartup(pod)
			if ready {
				pml.ready = ready
				pml.readyAt = readyAt
			}

			if err := pm.podUpdate(pml, pod, workspaceStack); err != nil {
				glog.Warningf("failed to process pod update %q: %v", pml.key(), err)
			}
		}()
	}
}

func (pm *podStartupLatencyDataMonitor) podUpdate(pml *podMilestones, pod *v1.Pod, workspaceStack string) error {
	pm.Lock()
	defer pm.Unlock()

	podKey := pml.key()
	if ml, ok := pm.pods[podKey]; ok {
		ml.merge(pml, pod)
		glog.V(4).Infof("podUpdate exists %q: %s", podKey, ml)
	} else {
		pm.pods[podKey] = pml
		glog.V(4).Infof("podUpdate new %q: %s", podKey, pml)
	}

	return pm.updateMetric(pm.pods[podKey], workspaceStack)
}

func (pm *podStartupLatencyDataMonitor) updateMetric(pml *podMilestones, workspaceStack string) error {
	key := pml.key()
	if !pm.initDone {
		glog.V(4).Infof("ignore metric for pod %q because of initial listing", key)
		return nil
	}

	if pml.ready && strings.Contains(pml.name, "workspace") {
		getWorkspacesMeasureTimes(pml, workspaceStack)
	}

	if pml.ready && !strings.Contains(pml.name, "workspace") && !strings.Contains(pml.name, "postgres") && !strings.Contains(pml.name, "keycloak") {
		containers := testContext.Containers{}
		for _, c := range pml.containers {
			containers.ContainerName = c.name
			containers.StartupLatency = pml.readyAt.Sub(c.pulledAt).Seconds()
			testContext.TestInstance.ChePods = append(testContext.TestInstance.ChePods, containers)
		}
	}
	return nil
}

func getWorkspacesMeasureTimes(pml *podMilestones, workspaceStack string) {
	if workspaceStack == "java-spring" {
		containersTimes := workspaceContainersTimes(pml)
		testContext.TestInstance.Workspaces.JavaMaven.Containers = containersTimes
	}
	if workspaceStack == "nodejs-stack" {
		containersTimes := workspaceContainersTimes(pml)
		testContext.TestInstance.Workspaces.NodeJS.Containers = containersTimes
	}
	if workspaceStack == "sample-stack" {
		containersTimes := workspaceContainersTimes(pml)
		testContext.TestInstance.Workspaces.Sample.Containers = containersTimes
	}
}

func workspaceContainersTimes(pml *podMilestones) []testContext.Containers {
	containerInstance := testContext.Containers{}
	var containers []testContext.Containers
	for _, c := range pml.containers {
		if c.started && !strings.Contains(c.name, "che-plugin-artifacts-broker") {
			containerInstance.ContainerName = c.name
			containerInstance.StartupLatency = pml.readyAt.Sub(c.startedAt).Seconds()
			containers = append(containers, containerInstance)
		}
	}
	return containers
}
