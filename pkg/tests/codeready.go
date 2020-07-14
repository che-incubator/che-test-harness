package tests

import (
	"bytes"
	"github.com/che-incubator/che-test-harness/pkg/client"
	logger2 "github.com/che-incubator/che-test-harness/pkg/controller/logger"
	"github.com/che-incubator/che-test-harness/pkg/monitors/metadata"
	"github.com/onsi/ginkgo"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Config Maps constants
	CodeReadyConfigMap     = "che"

	// Pod Names used to get info
	CodeReadyOperatorLabel = "olm.owner.kind=ClusterServiceVersion"
	PostgresLabel          = "component=postgres"
	KeycloakLabel          = "component=keycloak,app=codeready"
	DevFileLabel           = "component=devfile-registry"
	PluginRegistryLabel    = "component=plugin-registry"
	CodReadyServerLabel    = "component=codeready"

	//Custom Resource name to get info
	CRDName                = "checlusters.org.eclipse.che"

	// PVC name used for postgres-data
	PostgresPVCName        = "postgres-data"

	// Secret name used for add ca.crt
	chePostgresSecret      = "che-postgres-secret"
)

type PodInfo struct {
	Kind       string `json:"kind"`
	APIVersion string `json:"apiVersion"`
}

type CHE struct {
	token string `json:"access_token"`
}

var resource, _ = client.NewK8sClient()
var Logger = &logger2.Zap

// KubeDescribe is wrapper function for ginkgo describe. .
func KubeDescribe(text string, body func()) bool {
	return ginkgo.Describe("[Code Ready Workspaces Test Harness] "+text, body)
}

// GetConfigMap Return info about a specific config map
func GetConfigMap(configName string) (*v1.ConfigMap, error) {
	return resource.Kube().CoreV1().ConfigMaps(metadata.Namespace.Name).Get(configName, metav1.GetOptions{})
}

// GetSecret Return info about a specific secret into namespace
func GetSecret(secretName string) (*v1.Secret, error) {
	return resource.Kube().CoreV1().Secrets(metadata.Namespace.Name).Get(secretName, metav1.GetOptions{})
}

// GetPersistentVolumeClaims Return info about a specific PVC into namespace
func GetPersistentVolumeClaims(pvcName string) (*v1.PersistentVolumeClaim, error) {
	return resource.Kube().CoreV1().PersistentVolumeClaims(metadata.Namespace.Name).Get(pvcName, metav1.GetOptions{})
}

// GetServices return all services into namespace
func GetServices() (*v1.ServiceList, error) {
	return resource.Kube().CoreV1().Services(metadata.Namespace.Name).List(metav1.ListOptions{})
}

// GetPodByLabel return information about a specific pod
func GetPodByLabel(label string) (*v1.PodList, error) {
	return resource.Kube().CoreV1().Pods(metadata.Namespace.Name).List(metav1.ListOptions{LabelSelector: label})
}

// GetPodByLabel return information about a specific pod
func DeletePodByLabel(pod string) (error) {
	fg := metav1.DeletePropagationBackground
	deleteOptions := metav1.DeleteOptions{PropagationPolicy: &fg}
	return resource.Kube().CoreV1().Pods(metadata.Namespace.Name).Delete(pod, &deleteOptions)
}

// DescribePod set metadata and metrics about a specific pod
func DescribePod (pod *v1.PodList) (err error){
	var podInfo metadata.CodeReadyPods

	for _, v := range pod.Items {
		podInfo.Name = v.Name
		podInfo.Status = v.Status.Phase
		podInfo.Labels = v.Labels
		DescribePodLogs(v.Name)

		for _, val := range v.Spec.Containers  {
			podInfo.DockerImage = val.Image
		}
		a := append(metadata.Instance.CodeReadyPodsInfo, podInfo)

		metadata.Instance.CodeReadyPodsInfo = a
	}
	return err
}

// DescribePodLogs get all logs from a specific pod and write to a file
func DescribePodLogs(podName string)  {
	podLogOpts := v1.PodLogOptions{}

	Logger.Info("Obtaining logs from " + podName + " pod and writing them to file")
	req := resource.Kube().CoreV1().Pods(metadata.Namespace.Name).GetLogs(podName, &podLogOpts)
	podLogs, _ := req.Stream()

	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, _ = io.Copy(buf, podLogs)

	str := buf.Bytes()

	err := ioutil.WriteFile("/tmp/artifacts/che" + podName + ".log", str, 0644)
	if err != nil {
		Logger.Error("error writing logs to file")
	}
}
