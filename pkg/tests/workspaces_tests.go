package tests

import (
	"github.com/onsi/ginkgo"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/client"
	"gitlab.cee.redhat.com/codeready-workspaces/crw-osde2e/pkg/controller/workspaces"
	"go.uber.org/zap"
	"io/ioutil"
	"path/filepath"
	"time"
)

var _ =  KubeDescribe( "[Workspaces]" , func() {
	ginkgo.It("Start simple Workspace", func() {
		workspaceStack := "simple"
		httpClient, err := client.NewHttpClient()

		ctrl := workspaces.NewWorkspaceController(httpClient)

		fileLocation, err := filepath.Abs("samples/workspaces/workspace_sample.json")

		if err != nil {
			Logger.Panic("Failed to get workspace devfile from location ", zap.Error(err))
		}

		file, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			Logger.Panic("Failed to read workspace devfile ", zap.Error(err))
		}

		Logger.Info("Starting a new simple workspace")
		_ = ctrl.RunWorkspace(file, workspaceStack)
	})

	ginkgo.It("Start Python 3.7 Workspace", func() {
		workspaceStack := "python"
		httpClient, err := client.NewHttpClient()

		ctrl := workspaces.NewWorkspaceController(httpClient)

		fileLocation, err := filepath.Abs("samples/workspaces/workspace_python37.json")

		if err != nil {
			Logger.Panic("Failed to get workspace devFile from location ", zap.Error(err))
		}

		file, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			Logger.Panic("Failed to read workspace devfile ", zap.Error(err))
		}
		Logger.Info("Starting a new workspace with python 3.7")

		_ = ctrl.RunWorkspace(file, workspaceStack)
	})

	ginkgo.It("Start java Maven Workspace", func() {
		workspaceStack := "java-maven"
		httpClient, err := client.NewHttpClient()

		time.Sleep(10 *time.Second)
		ctrl := workspaces.NewWorkspaceController(httpClient)

		fileLocation, err := filepath.Abs("samples/workspaces/workspace_java_maven.json")

		if err != nil {
			Logger.Panic("Failed to get workspace devFile from location ", zap.Error(err))
		}

		file, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			Logger.Panic("Failed to read workspace devfile ", zap.Error(err))
		}
		Logger.Info("Starting a new workspace java-maven")
		_ = ctrl.RunWorkspace(file, workspaceStack)
	})
})