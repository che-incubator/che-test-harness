package performance

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/che-incubator/che-test-harness/internal/logger"
	"github.com/che-incubator/che-test-harness/pkg/client"
	"github.com/che-incubator/che-test-harness/pkg/deploy/workspaces"
	"github.com/onsi/ginkgo"
	"go.uber.org/zap"
)

var _ = ginkgo.Describe("[Workspaces]", func() {
	var Logger, err = logger.ZapLogger()
	if err != nil {
		panic("Failed to create zap logger")
	}

	ginkgo.It("Start java Maven Workspace", func() {
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
		workspace, _ := ctrl.RunWorkspace("java-spring", file)
		fmt.Println(workspace.ID)
	})

	ginkgo.It("Start NodeJS Workspace", func() {
		httpClient, err := client.NewHttpClient()

		ctrl := workspaces.NewWorkspaceController(httpClient)

		fileLocation, err := filepath.Abs("samples/workspaces/workspace_nodejs.json")

		if err != nil {
			Logger.Panic("Failed to get workspace devFile from location ", zap.Error(err))
		}

		file, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			Logger.Panic("Failed to read workspace devfile ", zap.Error(err))
		}
		Logger.Info("Starting a new Nodejs workspace")
		workspace, _ := ctrl.RunWorkspace("nodejs-stack", file)
		fmt.Println(workspace.ID)
	})

	ginkgo.It("Start Sample Workspace", func() {
		httpClient, err := client.NewHttpClient()

		ctrl := workspaces.NewWorkspaceController(httpClient)

		fileLocation, err := filepath.Abs("samples/workspaces/workspace_sample.json")

		if err != nil {
			Logger.Panic("Failed to get workspace devFile from location ", zap.Error(err))
		}

		file, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			Logger.Panic("Failed to read workspace devfile ", zap.Error(err))
		}
		Logger.Info("Starting a new sample workspace")
		workspace, _ := ctrl.RunWorkspace("sample-stack", file)
		fmt.Println(workspace.ID)
	})
})
