package tests

import (
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ =  KubeDescribe( "[PersistentVolumeClaims]" , func() {
	ginkgo.It("PVC `postgres-data` should be created", func() {
		Logger.Info("Check if PVC for postgres was created")
		secret, err := GetPersistentVolumeClaims(PostgresPVCName)

		if err != nil {
			Logger.Error("Error on getting info about pvc status")
		}

		Expect(secret).NotTo(BeNil())
		Expect(err).NotTo(HaveOccurred(), "failed to get pvc %v\n", PostgresPVCName)
	})
})
