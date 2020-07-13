package util

import (
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/client"
	log "gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/controller/logger"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/monitors/metadata"
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
	// Initialize Codeready Kubernetes client to create resources in a giving namespace
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		panic(err)
	}

	return k8sClient.Kube().CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
}

func OsdSetupNameSpace() bool {
	OsdNamespaces := []string{CrwNamespace, QeNamespace, NonOsdNamespace}
	Logger.Info("Start to detect namespace where CRW operator was deployed...")
	for _, namespace := range OsdNamespaces {
		_, err := GetNamespace(namespace)
		if err == nil {
			Logger.Info("Code Ready Workspaces detected on namespace: " + namespace)
			metadata.Namespace.Name = namespace

			return true
		}
	}

	Logger.Error("Error on start Code Ready Workspaces Test Harness. Please check provided namespace")
	return false
}
