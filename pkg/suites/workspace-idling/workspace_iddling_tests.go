package workspace_idling

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/che-incubator/che-test-harness/pkg/common/client"
	idling_defaults "github.com/che-incubator/che-test-harness/pkg/common/idling"
	"github.com/che-incubator/che-test-harness/pkg/common/logger"
	"github.com/che-incubator/che-test-harness/pkg/controller/workspaces"
	"github.com/onsi/ginkgo"
	"go.uber.org/zap"
)

var _ = ginkgo.Describe("[Workspace Idling]", func() {
	var Logger, err = logger.ZapLogger()
	if err != nil {
		panic("Failed to create zap logger")
	}
	httpClient, err := client.NewHttpClient()
	ctrl := workspaces.NewWorkspaceController(httpClient)
	var workspaceId string

	ginkgo.It("Create workspace", func() {
		fmt.Println("--- Creating workspace ---")
		fileLocation, err := filepath.Abs("samples/workspaces/workspace_java_maven.json")

		if err != nil {
			Logger.Panic("Failed to get workspace devFile from location ", zap.Error(err))
		}

		file, err := ioutil.ReadFile(fileLocation)
		if err != nil {
			Logger.Panic("Failed to read workspace devfile ", zap.Error(err))
		}

		workspaceId, err = ctrl.CreateWorkspace(idling_defaults.ConfigInstance.CheUrl, idling_defaults.ConfigInstance.KeycloakUrl, file)
		if err != nil {
			Logger.Panic("Failed to create workspace from devfile ", zap.Error(err))
		}
	})

	ginkgo.It("Start workspace", func() {
		fmt.Println("--- Starting workspace ---")
		if err := ctrl.StartWorkspace(idling_defaults.ConfigInstance.KeycloakUrl, idling_defaults.ConfigInstance.CheUrl, workspaceId); err != nil {
			Logger.Panic("Failed to start workspace", zap.Error(err))
		}

		if err := ctrl.WaitWorkspaceStatusViaApi(idling_defaults.ConfigInstance.KeycloakUrl, idling_defaults.ConfigInstance.CheUrl, workspaceId, "RUNNING", 600); err != nil {
			Logger.Panic("Waiting for status RUNNING for workspace id "+workspaceId+" failed ", zap.Error(err))
		}
	})

	ginkgo.It("Wait workspace to be idled", func() {
		fmt.Println("--- Waiting for workspace to be idled after " + strconv.Itoa(idling_defaults.ConfigInstance.IdlingTimeoutMinutes) + "+5 minutes ---")
		idlingTimeout := (idling_defaults.ConfigInstance.IdlingTimeoutMinutes + 5) * 60
		if err := ctrl.WaitWorkspaceStatusViaApi(idling_defaults.ConfigInstance.KeycloakUrl, idling_defaults.ConfigInstance.CheUrl, workspaceId, "STOPPED", idlingTimeout); err != nil {
			Logger.Panic("Waiting for workspace id "+workspaceId+" to be idled failed ", zap.Error(err))
		}
	})

	ginkgo.It("Stop workspace", func() {
		fmt.Println("--- Stopping worksapce ---")
		if err := ctrl.StopWorkspace(idling_defaults.ConfigInstance.KeycloakUrl, idling_defaults.ConfigInstance.CheUrl, workspaceId); err != nil {
			Logger.Panic("Failed to stop workspace", zap.Error(err))
		}

		if err := ctrl.WaitWorkspaceStatusViaApi(idling_defaults.ConfigInstance.KeycloakUrl, idling_defaults.ConfigInstance.CheUrl, workspaceId, "STOPPED", 180); err != nil {
			Logger.Panic("Waiting for status STOPPED for workspace id "+workspaceId+" failed ", zap.Error(err))
		}
	})

	ginkgo.It("Delete workspace", func() {
		fmt.Println("--- Deleting workspace ---")
		if err := ctrl.DeleteWorkspace(idling_defaults.ConfigInstance.KeycloakUrl, idling_defaults.ConfigInstance.CheUrl, workspaceId); err != nil {
			Logger.Panic("Failed to delete workspace", zap.Error(err))
		}

		if err := ctrl.WaitWorkspaceStatusViaApi(idling_defaults.ConfigInstance.KeycloakUrl, idling_defaults.ConfigInstance.CheUrl, workspaceId, "", 180); err != nil {
			Logger.Panic("Waiting for workspace with  id "+workspaceId+" to be removed failed ", zap.Error(err))
		}
	})

})
