package tests

import (
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ =  KubeDescribe( "[Services]", func() {
	ginkgo.It("Check if codeready services exist in cluster", func() {
		// TODO: Try to improve this checks with new err messages...
		Logger.Info("Checking if all services for Code Ready Workspaces")
		services, err := GetServices()

		Expect(services).NotTo(BeNil())

		confmap := map[string]string{}
		for _ ,v:= range services.Items {
			confmap[v.Name]= v.Name
		}

		Expect(confmap["che-host"]).NotTo(BeEmpty())
		Expect(confmap["plugin-registry"]).NotTo(BeEmpty())
		Expect(confmap["postgres"]).NotTo(BeEmpty())
		Expect(confmap["devfile-registry"]).NotTo(BeEmpty())

		Expect(err).NotTo(HaveOccurred())
	})
})

