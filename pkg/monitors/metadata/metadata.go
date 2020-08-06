package metadata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"os"
)

// metadata houses metadata to be written out to the additional-metadata.json
type metadata struct {
	// Whether the CRD was found. Typically Spyglass seems to have issues displaying non-strings, so
	// this will be written out as a string despite the native JSON boolean type.
	FoundCRD bool `json:"found-crd,string"`

	ClusterTimeUp float64 `json:"cluster-time-up, int"`

	ChePodsInfo []ChePodsInfo `json:"pods-info,string"`

	ChePodTime PodTimes `json:"pods-up-times, int"`

	WorkspacesMeasureTime WorkspacesMeasureTime `json:"workspaces-up-time, int"`

	//Returns true or false depending if che server is UP
	CheServerIsUp bool `json:"che_apiserver_is_up, bool"`
}

type PodTimes struct {
	PostgresUpTime  float64 `json:"postgres-up-time, float64"`
	KeycloakUpTime  float64 `json:"keycloak-up-time, float64"`
	DevFileUpTime   float64 `json:"devfile-up-time, float64"`
	PluginRegUpTime float64 `json:"plugins-up-time, float64"`
	ServerUpTime float64 `json:"server-up-time, float64"`
}

type WorkspacesMeasureTime struct {
	SimpleWorkspace float64 `json:"simple_workspace, float64"`
	JavaMavenWorkspace float64 `json:"java_maven_workspace, float64"`
}

type ChePodsInfo struct {
	Name        string             `json:"name, string"`
	DockerImage string             `json:"docker_image, string"`
	Status      v1.PodPhase        `json:"status, string"`
	Labels      map[string]string  `json:"labels, string"`
}

type CHE_NAMESPACE struct {
	Name string
	UP bool
}

var Namespace = CHE_NAMESPACE{}
var Instance  = metadata{}

// WriteToJSON will marshall the metadata struct a	nd write it into the given file.
func (m *metadata) WriteToJSON(outputFilename string) (err error) {
	var data []byte
	if data, err = json.Marshal(m); err != nil {
		return err
	}

	if err = ioutil.WriteFile(outputFilename, data, os.FileMode(0644)); err != nil {
		return err
	}

	return nil
}

// WriteWorkspaceMeasureTimesToMetadata get timers condition from a given pod
func WriteWorkspaceMeasureTimesToMetadata(pod *v1.Pod, workspaceStack string) (err error) {
	if pod.Status.Phase == v1.PodRunning {
		timeDiff := GetMeasureTime(pod)

		switch workspaceStack {
		case "simple":
			Instance.WorkspacesMeasureTime.SimpleWorkspace = timeDiff
		case "java-maven" :
			Instance.WorkspacesMeasureTime.JavaMavenWorkspace = timeDiff
		}

	} else {
		return fmt.Errorf("Error on check workspace Pod Timers Pod didn't start.")
	}
	
	return 
}

// Get Measure time for Workspaces
// !TODO get with more specific events
func GetMeasureTime(pod *v1.Pod) (time float64) {
	for _, p := range pod.Status.Conditions   {
		if p.Type == "ContainersReady" {
			startupTime := p.LastTransitionTime.Time.Sub(pod.Status.StartTime.Time).Seconds()

			return startupTime
		}
	}
	return time
}
