package tests

import (
	"github.com/che-incubator/che-test-harness/pkg/monitors/metadata"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var _ = KubeDescribe( "[Custom Resources]", func() {
	ginkgo.It("Check if CRD already exist in Cluster", func() {
		Logger.Info("Checking if CRD for Che exist in cluster")
		// Move this client
		cfg, err := config.GetConfig()
		apiextensions, err := clientset.NewForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())
		// Make sure the CRD exist in cluster
		_, err = apiextensions.ApiextensionsV1beta1().CustomResourceDefinitions().Get(CRDName, metav1.GetOptions{})
		if err != nil {
			metadata.Instance.FoundCRD = false
		} else {
			metadata.Instance.FoundCRD = true
		}

		Expect(err).NotTo(HaveOccurred())
	})
})
