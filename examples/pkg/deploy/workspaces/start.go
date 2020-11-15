package workspaces

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/che-incubator/che-test-harness/pkg/client"
	"github.com/che-incubator/che-test-harness/pkg/deploy"
	testContext "github.com/che-incubator/che-test-harness/pkg/deploy/context"
	"go.uber.org/zap"
	"net/http"
)

const (
	ClientID = "che-public"
)

// RunWorkspace create a new performance from a given devfile and call an method to get measure time for performance after a performance pod is up and ready
func (w *Controller) RunWorkspace(workspaceStackName string,workspaceDefinition []byte) (workspace *Workspace, err error) {
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		w.Logger.Panic("Failed to create kubernetes client.", zap.Error(err))
	}

	ctrl := deploy.NewTestHarnessController(k8sClient)
	resource, err := ctrl.GetCustomResource()
	if err != nil {
		w.Logger.Panic("Failed to get Custom Resource.", zap.Error(err))
	}

	workspace, err = w.CreateWorkspace(resource.Status.CheURL, resource.Status.KeycloakURL+testContext.KeycloakTokenEndpoint, workspaceDefinition, workspaceStackName)
	if err != nil {
		w.Logger.Panic("Error on create performance.", zap.Error(err))
	}

	err = w.DeleteWorkspace(resource.Status.KeycloakURL+testContext.KeycloakTokenEndpoint, resource.Status.CheURL, workspace)

	return workspace, err
}

// CreateWorkspace create an performance using token and given devFile
func (w *Controller) CreateWorkspace(cheURL string, keycloakUrl string, workspaceDefinition []byte, workspaceStackName string) (workspace *Workspace, err error) {
	var keycloakAuth *KeycloakAuth

	if keycloakAuth, err = w.KeycloakToken(keycloakUrl); err != nil {
		w.Logger.Error("Failed to get user token ", zap.Error(err))
	}

	request, _ := http.NewRequest("POST", cheURL+"/api/workspace/devfile", bytes.NewBuffer(workspaceDefinition))
	request.Header.Add("Authorization", "Bearer "+keycloakAuth.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	response, err := w.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.Status == "409" {
		return nil, errors.New("workspace already exist")
	}

	err = json.NewDecoder(response.Body).Decode(&workspace)
	startWorkspace, err := w.startWorkspace(cheURL, keycloakAuth.AccessToken, workspace.ID)
	if err != nil && !startWorkspace {
		return nil, errors.New("error sending request to start workspace")
	}

	w.Logger.Info("Waiting 5 minutes to start a workspace")
	statusWorkspace, err := w.statusWorkspace(cheURL, keycloakUrl, workspace, workspaceStackName, testContext.WorkspaceRunningStatus)
	if !statusWorkspace {
		return nil, err
	}
	return workspace, err
}

func (w *Controller) startWorkspace(cheUrl string, token string, workspaceID string) (boolean bool, err error) {
	request, err := http.NewRequest("POST", cheUrl+"/api/workspace/"+workspaceID+"/runtime", nil)
	if err != nil {
		return
	}

	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")

	_, err = w.httpClient.Do(request)

	if err != nil {
		return false, err
	}

	return true, err
}
