package che

import (
	"github.com/che-incubator/che-test-harness/pkg/common/client"
	"github.com/che-incubator/che-test-harness/pkg/common/logger"
	"github.com/che-incubator/che-test-harness/pkg/controller/workspaces"
	"github.com/onsi/ginkgo"
	"go.uber.org/zap"
	"io/ioutil"
	"path/filepath"
)

var _ =  ginkgo.Describe( "[Workspaces]" , func() {
	var Logger, err = logger.ZapLogger()
	if err != nil {
		panic("Failed to create zap logger")
	}

	ginkgo.It("Start java Maven Workspace", func() {
		workspaceStack := "java-maven"
		httpClient, err := client.NewHttpClient()

		ctrl := workspaces.NewWorkspaceController(httpClient)

		fileLocation, err := filepath.Abs("samples/workspaces/workspace_java_maven.json")

		if err != nil {
			Logger.Panic("Failed to get workspace devFile from location ", zap.Error(err))
		}

		file, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			Logger.Panic("Failed to read workspace devfile ", zap.Error(err))
		}
		Logger.Info("Starting a new Java Maven workspace")
		_ = ctrl.RunWorkspace(file, workspaceStack)
	})
})
