package monitors

import (
	"context"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/cmd/operator_osd/config"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/controller/logger"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/monitors/metadata"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"

	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/pager"

	kubeletevents "k8s.io/kubernetes/pkg/kubelet/events"
)

// PodStartupMilestones keeps all milestone timestamps from Pod creation.
type PodStartupMilestones struct {
	created          time.Time
	startedPulling   time.Time
	finishedPulling  time.Time
	observedRunning  time.Time
	seenPulled       int
	needPulled       int
}

// PodStartupLatencyDataMonitor monitors pod startup latency and exposes prometheus metric.
type PodStartupLatencyDataMonitor struct {
	sync.Mutex
	kubeClient     clientset.Interface
	PodStartupData map[string]PodStartupMilestones
	logger *zap.Logger
}

// IsComplete returns true is data is complete (ready to be included in the metric) and if it haven't been included in the metric yet.
func (data *PodStartupMilestones) IsComplete() bool {
	return !data.created.IsZero() && !data.startedPulling.IsZero() && !data.finishedPulling.IsZero() && !data.observedRunning.IsZero() && data.seenPulled == data.needPulled //&& !data.accountedFor
}

// NewPodStartupLatencyDataMonitor creates a new PodStartupLatencyDataMonitor from a given client.
func NewPodStartupDataMonitor(c clientset.Interface) *PodStartupLatencyDataMonitor {
	l, err := logger.ZapLogger()

	if err != nil {
		panic(err)
	}
	return &PodStartupLatencyDataMonitor{
		kubeClient:     c,
		PodStartupData: map[string]PodStartupMilestones{},
		logger: l,
	}
}

// Descibe all pod events in given namespace for Code Ready Workspaces
func (pm *PodStartupLatencyDataMonitor) DescribeEvents(stopCh chan struct{}) error {
	_, controller := cache.NewInformer(&cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			pg := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
				return pm.kubeClient.CoreV1().Pods(metadata.Namespace.Name).List(opts)
			}))
			return pg.List(context.Background(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			return pm.kubeClient.CoreV1().Pods(metadata.Namespace.Name).Watch(options)
		},
	}, new(v1.Pod), 0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {},
			DeleteFunc: func(obj interface{}) {},
			UpdateFunc: func(oldObj, newObj interface{}) {
				p := newObj.(*v1.Pod)
				pm.handlePodUpdate(p)
			},
		})
	go controller.Run(stopCh)

	eventSelector := fields.Set{
		"involvedObject.kind": "Pod",
		"source": "kubelet",
	}.AsSelector().String()

	_, eventcontroller := cache.NewInformer(&cache.ListWatch{
		ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
			options.FieldSelector = eventSelector
			pg := pager.New(pager.SimplePageFunc(func(opts metav1.ListOptions) (runtime.Object, error) {
				return pm.kubeClient.CoreV1().Events(metadata.Namespace.Name).List(opts)
			}))
			return pg.List(context.Background(), options)
		},
		WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
			options.FieldSelector = eventSelector
			return pm.kubeClient.CoreV1().Events(metadata.Namespace.Name).Watch(options)
		},
	}, new(v1.Event), 0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				e := obj.(*v1.Event)
				pm.handleEventUpdate(e)
			},
			DeleteFunc: func(obj interface{}) {},
		})
	go eventcontroller.Run(stopCh)
	return nil
}

func (pm *PodStartupLatencyDataMonitor) handlePodUpdate(p *v1.Pod) {
	if isReady, create, running := checkPodAndGetStartup(p); isReady {
		go pm.podRunningTime(getPodKey(p), create, running, len(p.Spec.Containers))
	}
}

func (pm *PodStartupLatencyDataMonitor) handlePullingImageEvent(key string, e *v1.Event) {
	pm.Lock()
	defer pm.Unlock()

	ok := false
	data := PodStartupMilestones{}
	if data, ok = pm.PodStartupData[key]; !ok {
		data.finishedPulling = time.Unix(0, 0)
		data.needPulled = -1
	}
	if data.startedPulling.IsZero() || data.startedPulling.After(e.FirstTimestamp.Time) {
		data.startedPulling = e.FirstTimestamp.Time
	}

	pm.updateMetric(key, &data)
	pm.PodStartupData[key] = data
}

func (pm *PodStartupLatencyDataMonitor) handlePulledImageEvent(key string, e *v1.Event) {
	pm.Lock()
	defer pm.Unlock()

	ok := false
	data := PodStartupMilestones{}
	if data, ok = pm.PodStartupData[key]; ok {
		data.startedPulling = time.Unix(0, 0)
		data.needPulled = -1
	}
	// Check if image is already pulled in machine.
	if data.finishedPulling.IsZero() || data.finishedPulling.Before(e.FirstTimestamp.Time) {
		data.finishedPulling = e.FirstTimestamp.Time
	}
	if strings.Contains(e.Message, "already present on machine") {
		data.startedPulling = e.FirstTimestamp.Time
	}
	data.seenPulled++

	pm.updateMetric(key, &data)
	pm.PodStartupData[key] = data
}

func (pm *PodStartupLatencyDataMonitor) handleEventUpdate(e *v1.Event) {
	key := getPodKeyFromReference(&e.InvolvedObject)
	switch e.Reason {
	case kubeletevents.PullingImage:
		go pm.handlePullingImageEvent(key, e)
	case kubeletevents.PulledImage:
		go pm.handlePulledImageEvent(key, e)

	default:
		return
	}
}

// Detect if current event for given pod is complete and save measure times
func (pm *PodStartupLatencyDataMonitor) updateMetric(key string, data *PodStartupMilestones) {
	if data.IsComplete() {
		finishedPulling := data.finishedPulling
		observeRunning := data.observedRunning

		startupTime := observeRunning.Sub(finishedPulling).Seconds()
		pm.logger.Info("Successfully created pod " + key, zap.Time("createdAt",data.created))

		pm.logger.Info("Pulled successfully container image for pod: " + key, zap.Time("pulled at",data.finishedPulling))
		pm.logger.Info("Pod " + key + " successfully succeeded:", zap.Time("readyAt", data.observedRunning))

		// TODO: Find a better way to identify pods running
		if strings.Contains(key, "postgres") {
			metadata.Instance.CRWPodTime.PostgresUpTime= startupTime
		}

		if strings.Contains(key, "keycloak") {
			metadata.Instance.CRWPodTime.KeycloakUpTime= startupTime
		}

		if strings.Contains(key, "plugin-registry") {
			metadata.Instance.CRWPodTime.PluginRegUpTime= startupTime
		}

		if strings.Contains(key, "devfile-registry") {
			metadata.Instance.CRWPodTime.DevFileUpTime= startupTime
		}

		if strings.Contains(key, config.TestHarnessConfig.Flavor) && ! strings.Contains(key, "codeready-operator") {
			metadata.Instance.CRWPodTime.ServerUpTime= startupTime
		}
	}
}

func (pm *PodStartupLatencyDataMonitor) podRunningTime(podKey string, createTime time.Time, runningTime time.Time, needPulled int) {
	pm.Lock()
	defer pm.Unlock()

	var data PodStartupMilestones
	var ok bool
	if data, ok = pm.PodStartupData[podKey]; !ok {
		// Necessary to work anywhere except UTC time zone...
		data.startedPulling = time.Unix(0, 0)
		data.finishedPulling = time.Unix(0, 0)
	}
	data.created = createTime
	data.observedRunning = runningTime
	data.needPulled = needPulled

	pm.updateMetric(podKey, &data)
	pm.PodStartupData[podKey] = data
}
