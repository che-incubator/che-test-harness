package tests

import (
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _  = KubeDescribe( "[ConfigMaps]" ,func() {
	ginkgo.It("Config map `che` should exist", func() {
		Logger.Info("Checking `che` config map integrity")
		che, err := GetConfigMap(CodeReadyConfigMap)

		Expect(che).NotTo(BeNil())

		if err != nil {
			Logger.Error("Error on verify `che` config map")
		}

		Expect(err).NotTo(HaveOccurred())
	})
})
