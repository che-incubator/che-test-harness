package workspaces

import (
	"bytes"
	"encoding/json"
	"github.com/che-incubator/che-test-harness/pkg/common/client"
	"github.com/che-incubator/che-test-harness/pkg/common/logger"
	"github.com/che-incubator/che-test-harness/pkg/controller"
	"github.com/che-incubator/che-test-harness/pkg/monitors/metadata"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
)

// WorkspacesController useful to add logger and http client.
type WorkspacesController struct {
	httpClient *http.Client
	Logger     logger.Log
}

// NewWorkspaceController creates a new WorkspacesController from a given client.
func NewWorkspaceController(c *http.Client) *WorkspacesController {
	return &WorkspacesController{
		httpClient: c,
		Logger:     logger.Zap,
	}
}

// RunWorkspace create a new workspace from a given devfile and call an method to get measure time for workspace after a workspace pod is up and ready
func (w *WorkspacesController) RunWorkspace(workspaceDefinition []byte, workspaceStack string) (workspaceID string) {
	k8sClient, err := client.NewK8sClient()

	if err != nil {
		w.Logger.Panic("Failed to create kubernetes client.", zap.Error(err))
	}

	ctrl := controller.NewTestHarnessController(k8sClient)

	resource, err := ctrl.GetCustomResource()
	if err != nil {
		w.Logger.Panic("Failed to get Custom Resource.", zap.Error(err))
	}

	keycloakTokenUrl := resource.Status.KeycloakURL
	cheURL := resource.Status.CheURL

	accessToken, err := w.KeycloakToken(keycloakTokenUrl + "/auth/realms/che/protocol/openid-connect/token/")
	if err != nil {
		w.Logger.Panic("Error on retrieving token from keycloak.", zap.Error(err))
	}

	workspaceID, err = w.CreateWorkspace(cheURL, accessToken, workspaceDefinition)
	if err != nil {
		w.Logger.Panic("Error on create workspace.", zap.Error(err))
	}

	if err := w.StartWorkspace(accessToken, cheURL, workspaceID); err != nil {
		w.Logger.Panic("Failed to start workspace", zap.Error(err))
	}

	workspaceLabel := "che.workspace_id=" + workspaceID

	if _, err := ctrl.WatchPodStartup(metadata.Namespace.Name, workspaceLabel, workspaceStack); err != nil {
		w.Logger.Panic("Failed to start workspace", zap.Error(err))
	}

	if err := w.StopWorkspace(accessToken, cheURL, workspaceID); err != nil {
		w.Logger.Panic("Failed to delete workspace", zap.Error(err))
	}

	if err := w.DeleteWorkspace(accessToken, cheURL, workspaceID); err != nil {
		w.Logger.Panic("Failed to delete workspace", zap.Error(err))
	}

	return workspaceID
}

// KeycloakToken return a JWT from keycloak
func (w *WorkspacesController) KeycloakToken(keycloakTokenUrl string) (token string, err error) {
	var result map[string]interface{}
	cheFlavor := "che"

	data := url.Values{}

	data.Set("client_id", cheFlavor + "-public")
	data.Set("username", "admin")
	data.Set("password", "admin")
	data.Set("grant_type", "password")

	r, err := http.NewRequest("POST", keycloakTokenUrl, strings.NewReader(data.Encode()))

	if err != nil {
		w.Logger.Panic("Failed to get token from keycloak", zap.Error(err))
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := w.httpClient.Do(r)

	err = json.NewDecoder(response.Body).Decode(&result)

	token = result["access_token"].(string)

	return token, err
}

// CreateWorkspace create an workspace using token and given devFile
func (w *WorkspacesController) CreateWorkspace(cheURL string, token string, workspaceDefinition []byte) (workspaceID string, err error) {
	var workspace map[string]interface{}

	request, _ := http.NewRequest("POST", cheURL+"/api/workspace/devfile?infrastructure-namespace=" + metadata.Namespace.Name + "&namespace=admin", bytes.NewBuffer(workspaceDefinition))
	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")

	response, err := w.httpClient.Do(request)

	if err != nil && response.Status != "201" {
		w.Logger.Panic("Error on create workspace...", zap.Error(err))
	}

	_ = json.NewDecoder(response.Body).Decode(&workspace)

	workspaceID = workspace["id"].(string)

	if len(workspaceID) == 0 {
		w.Logger.Panic("Workspace ID is empty.The tests will fail.")
	}

	return workspaceID, err
}

// StartWorkspace start a new workspace from a given workspace_id
func (w *WorkspacesController) StartWorkspace(token string, cheURL string, workspaceID string) (err error) {
	request, err := http.NewRequest("POST", cheURL + "/api/workspace/" + workspaceID + "/runtime", nil)

	if err != nil {
		w.Logger.Error("Failed to start Workspace", zap.Error(err))
	}

	request.Header.Add("Authorization", "Bearer " + token)
	request.Header.Add("Content-Type", "application/json")

	res, err := w.httpClient.Do(request)

	if res.Status != "200" && err != nil {
		return err
	}

	return err
}

// StopWorkspace stop a workspace from a given workspace_id
func (w *WorkspacesController) StopWorkspace(token string, cheURL string, workspaceID string) (err error) {
	request, err := http.NewRequest("DELETE", cheURL + "/api/workspace/" + workspaceID + "/runtime", nil)

	if err != nil {
		w.Logger.Error("Failed to stop workspace", zap.Error(err))
	}

	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")

	_ , err = w.httpClient.Do(request)

	return err
}

// DeleteWorkspace delete a workspace from a given workspace_id
func (w *WorkspacesController) DeleteWorkspace(token string, cheURL string, workspaceID string) (err error) {
	request, err := http.NewRequest("DELETE", cheURL + "/api/workspace/" + workspaceID, nil)

	if err != nil {
		w.Logger.Error("Failed to delete workspace", zap.Error(err))
	}

	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")

	_ , err = w.httpClient.Do(request)

	return err
}
