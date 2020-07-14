package util

import (
	"github.com/che-incubator/che-test-harness/pkg/client"
	log "github.com/che-incubator/che-test-harness/pkg/controller/logger"
	"github.com/che-incubator/che-test-harness/pkg/monitors/metadata"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

var (
	QeNamespace     = "codeready-workspaces-operator-qe"
	CrwNamespace    = "codeready-workspaces-operator"
	NonOsdNamespace = os.Getenv("CODEREADY_NAMESPACE")
)

var Logger = &log.Zap

func GetNamespace(namespace string) (*v1.Namespace, error) {
	// Initialize Kubernetes client to create resources in a giving namespace
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		panic(err)
	}

	return k8sClient.Kube().CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
}

func OsdSetupNameSpace() bool {
	OsdNamespaces := []string{CrwNamespace, QeNamespace, NonOsdNamespace}
	Logger.Info("Start to detect namespace where Che operator was deployed...")
	for _, namespace := range OsdNamespaces {
		_, err := GetNamespace(namespace)
		if err == nil {
			Logger.Info("Eclipse Che Operator detected on namespace: " + namespace)
			metadata.Namespace.Name = namespace

			return true
		}
	}

	Logger.Error("Error on start Eclipse Che Test Harness. Please check provided namespace")
	return false
}
