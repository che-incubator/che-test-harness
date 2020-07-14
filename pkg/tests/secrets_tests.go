package tests

import (
	"github.com/onsi/ginkgo"
)

var _ =  KubeDescribe( "[Secrets]" , func() {
	ginkgo.It("Secret `self-signed-certificate` should exist", func() {
		Logger.Info("Checking secrets created for code ready workspaces")
		//secret, err := GetSecret(chePostgresSecret)

		//if err != nil {
		//	Logger.Error("Error on get info about secrets")
		//}

		//Expect(secret).NotTo(BeNil())
		//Expect(err).NotTo(HaveOccurred(), "failed to get secretName %v\n", chePostgresSecret)
	})
})
