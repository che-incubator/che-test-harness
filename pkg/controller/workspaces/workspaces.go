package workspaces

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/che-incubator/che-test-harness/pkg/common/client"
	"github.com/che-incubator/che-test-harness/pkg/common/logger"
	"github.com/che-incubator/che-test-harness/pkg/controller"
	"github.com/che-incubator/che-test-harness/pkg/monitors/metadata"
	"go.uber.org/zap"

	. "github.com/onsi/gomega"
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

	keycloakTokenUrl := resource.Status.KeycloakURL + "/auth/realms/che/protocol/openid-connect/token/"
	cheURL := resource.Status.CheURL

	workspaceID, err = w.CreateWorkspace(cheURL, keycloakTokenUrl, workspaceDefinition)
	if err != nil {
		w.Logger.Panic("Error on create workspace.", zap.Error(err))
	}

	if err := w.StartWorkspace(keycloakTokenUrl, cheURL, workspaceID); err != nil {
		w.Logger.Panic("Failed to start workspace", zap.Error(err))
	}

	workspaceLabel := "che.workspace_id=" + workspaceID

	if _, err := ctrl.WatchPodStartup(metadata.Namespace.Name, workspaceLabel, workspaceStack); err != nil {
		w.Logger.Panic("Failed to start workspace", zap.Error(err))
	}

	if err := w.StopWorkspace(keycloakTokenUrl, cheURL, workspaceID); err != nil {
		w.Logger.Panic("Failed to delete workspace", zap.Error(err))
	}

	if err := w.DeleteWorkspace(keycloakTokenUrl, cheURL, workspaceID); err != nil {
		w.Logger.Panic("Failed to delete workspace", zap.Error(err))
	}

	return workspaceID
}

// KeycloakToken return a JWT from keycloak
func (w *WorkspacesController) KeycloakToken(keycloakTokenUrl string) (token string, err error) {
	var result map[string]interface{}
	cheFlavor := "che"

	data := url.Values{}

	data.Set("client_id", cheFlavor+"-public")
	data.Set("username", metadata.User.Username)
	data.Set("password", metadata.User.Password)
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
func (w *WorkspacesController) CreateWorkspace(cheURL string, keycloakUrl string, workspaceDefinition []byte) (workspaceID string, err error) {
	var workspace map[string]interface{}
	var request *http.Request

	if metadata.Namespace.Name == "" {
		request, _ = http.NewRequest("POST", cheURL+"/api/workspace/devfile?", bytes.NewBuffer(workspaceDefinition))
	} else {
		request, _ = http.NewRequest("POST", cheURL+"/api/workspace/devfile?infrastructure-namespace="+metadata.Namespace.Name+"&namespace=admin", bytes.NewBuffer(workspaceDefinition))
	}

	var token string
	if token, err = w.KeycloakToken(keycloakUrl); err != nil {
		w.Logger.Error("Failed to get user token ", zap.Error(err))
	}

	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")

	response, err := w.httpClient.Do(request)

	if err != nil && response.Status != "201" {
		w.Logger.Panic("Error on create workspace...", zap.Error(err))
	}

	_ = json.NewDecoder(response.Body).Decode(&workspace)

	if response.StatusCode == 409 {
		w.Logger.Panic("Can not create workspace, workspace probably exists.")
	}

	workspaceID = workspace["id"].(string)

	if len(workspaceID) == 0 {
		w.Logger.Panic("Workspace ID is empty.The tests will fail.")
	}

	return workspaceID, err
}

// StartWorkspace start a new workspace from a given workspace_id
func (w *WorkspacesController) StartWorkspace(keycloakUrl string, cheURL string, workspaceID string) (err error) {
	request, err := http.NewRequest("POST", cheURL+"/api/workspace/"+workspaceID+"/runtime", nil)

	if err != nil {
		w.Logger.Error("Failed to start Workspace", zap.Error(err))
	}

	var token string
	if token, err = w.KeycloakToken(keycloakUrl); err != nil {
		w.Logger.Error("Failed to get user token ", zap.Error(err))
	}

	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")

	res, err := w.httpClient.Do(request)

	if res.Status != "200" && err != nil {
		return err
	}

	return err
}

// StopWorkspace stop a workspace from a given workspace_id
func (w *WorkspacesController) StopWorkspace(keycloakUrl string, cheURL string, workspaceID string) (err error) {
	request, err := http.NewRequest("DELETE", cheURL+"/api/workspace/"+workspaceID+"/runtime", nil)

	if err != nil {
		w.Logger.Error("Failed to stop workspace", zap.Error(err))
	}

	var token string
	if token, err = w.KeycloakToken(keycloakUrl); err != nil {
		w.Logger.Error("Failed to get user token ", zap.Error(err))
	}

	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")

	_, err = w.httpClient.Do(request)

	return err
}

// DeleteWorkspace delete a workspace from a given workspace_id
func (w *WorkspacesController) DeleteWorkspace(keycloakUrl string, cheURL string, workspaceID string) (err error) {
	request, err := http.NewRequest("DELETE", cheURL+"/api/workspace/"+workspaceID, nil)

	if err != nil {
		w.Logger.Error("Failed to delete workspace", zap.Error(err))
	}

	var token string
	if token, err = w.KeycloakToken(keycloakUrl); err != nil {
		w.Logger.Error("Failed to get user token ", zap.Error(err))
	}

	request.Header.Add("Authorization", "Bearer "+token)
	request.Header.Add("Content-Type", "application/json")

	_, err = w.httpClient.Do(request)

	return err
}

// creating get requested with always fresh user token
func (w *WorkspacesController) getRequest(keycloakUrl string, cheUrl string, workspaceID string) *http.Request {
	request, err := http.NewRequest("GET", cheUrl+"/api/workspace/"+workspaceID, nil)

	if err != nil {
		w.Logger.Error("Failed to create request for obtaining workspace status ", zap.Error(err))
	}

	var token string
	if token, err = w.KeycloakToken(keycloakUrl); err != nil {
		w.Logger.Error("Failed to get user token ", zap.Error(err))
	}

	request.Header.Add("Authorization", "Bearer "+token)
	return request
}

// WaitWorkspaceStatusViaApi waits for workspace to have desired status in specified timeout
func (w *WorkspacesController) WaitWorkspaceStatusViaApi(keycloakUrl string, cheURL string, workspaceID string, desiredStatus string, timeoutInSeconds int) (err error) {
	var result map[string]interface{}

	startTime, endTime := time.Now(), time.Now()

	request := w.getRequest(keycloakUrl, cheURL, workspaceID)

	status := "undefined"
	response, err := w.httpClient.Do(request)

	if err != nil {
		w.Logger.Error("Failed to obtain workspace status ", zap.Error(err))
	}

	err = json.NewDecoder(response.Body).Decode(&result)

	timeouted := false
	fmt.Println("Waiting for workspace to be " + desiredStatus + ". Tick is set to 5 sec.")
	for status != desiredStatus {
		request := w.getRequest(keycloakUrl, cheURL, workspaceID)
		response, err := w.httpClient.Do(request)

		if err != nil {
			w.Logger.Error("Failed to obtain workspace status ", zap.Error(err))
		}

		err = json.NewDecoder(response.Body).Decode(&result)

		if result["status"] == nil {
			if desiredStatus == "" {
				break
			} else {
				statusError := errors.New("Can not obtain workspace status. Status is empty.")
				w.Logger.Error("Can not obtain workspace status ", zap.Error(statusError))
				Expect(statusError).NotTo(HaveOccurred())
			}
		}

		status = result["status"].(string)
		endTime = time.Now()
		fmt.Println("Status: ", status, "Wanted: ", desiredStatus, " Time taken: ", endTime.Sub(startTime))

		if endTime.Sub(startTime) > time.Duration(timeoutInSeconds)*(time.Second) {
			fmt.Println("Timeouting from after: ", endTime.Sub(startTime), " from timeout: ", timeoutInSeconds)
			timeouted = true
			break
		}

		time.Sleep(5 * time.Second)
	}

	if timeouted {
		timeoutError := errors.New("Waiting for workspace to get to the state " + desiredStatus + " failed, workspace is " + status)
		w.Logger.Error("Waiting for workspace to change status timeouted ", zap.Error(timeoutError))
		Expect(timeoutError).NotTo(HaveOccurred())
	} else {
		fmt.Println("Workspace become ", status, " after ", endTime.Sub(startTime))
	}

	return err
}
