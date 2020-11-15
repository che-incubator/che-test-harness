package monitors

import (
	"fmt"
	"strings"
	"time"

	"k8s.io/api/core/v1"
)

func fetchContainerName(objRef v1.ObjectReference) (string, error) {
	name := objRef.FieldPath
	if strings.HasPrefix(name, "spec.containers{") {
		name = strings.TrimPrefix(name, "spec.containers{")
		name = strings.TrimSuffix(name, "}")
		return name, nil
	}

	if strings.HasPrefix(name, "spec.initContainers{") {
		name = strings.TrimPrefix(name, "spec.initContainers{")
		name = strings.TrimSuffix(name, "}")
		return name, nil
	}

	return "", fmt.Errorf("unknown format: %s", name)
}

func getPodKey(namespace, name, uid string) string {
	return namespace + "/" + name + "/" + uid
}

func getPodKeyFromReference(objRef v1.ObjectReference) string {
	return getPodKey(objRef.Namespace, objRef.Name, string(objRef.UID))
}

func podReady(pod *v1.Pod) bool {
	for _, cond := range pod.Status.Conditions {
		if cond.Type == v1.PodReady && cond.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

func podRunningReady(p *v1.Pod) (bool, error) {
	// Check the phase is running.
	if p.Status.Phase != v1.PodRunning {
		return false, fmt.Errorf("want pod '%s' on '%s' to be '%v' but was '%v'",
			p.ObjectMeta.Name, p.Spec.NodeName, v1.PodRunning, p.Status.Phase)
	}
	// Check the ready condition is true.
	if !podReady(p) {
		return false, fmt.Errorf("pod '%s' on '%s' didn't have condition {%v %v}; conditions: %v",
			p.ObjectMeta.Name, p.Spec.NodeName, v1.PodReady, v1.ConditionTrue, p.Status.Conditions)
	}
	return true, nil
}

// PodRunningReadyOrSucceeded checks whether pod p's phase is running and it has a ready
// condition of status true or wheather the Pod already succeeded.
func PodRunningReadyOrSucceeded(p *v1.Pod) (bool, error) {
	// Check if the phase is succeeded.
	if p.Status.Phase == v1.PodSucceeded {
		return true, nil
	}
	return podRunningReady(p)
}

func checkPodAndGetStartup(p *v1.Pod) (bool, time.Time) {
	if ok, _ := PodRunningReadyOrSucceeded(p); !ok {
		return false, time.Unix(0, 0)
	}
	return true, time.Now()
}
