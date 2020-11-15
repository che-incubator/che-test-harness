package workspaces

import (
	testContext "github.com/che-incubator/che-test-harness/pkg/deploy/context"
	"go.uber.org/zap"
	"net/http"
)

// DeleteWorkspace delete a workspace from a given workspace_id
func (w *Controller) DeleteWorkspace(keycloakUrl string, cheURL string, workspace *Workspace) (err error) {
	var keycloakAuth *KeycloakAuth

	if err = w.stopWorkspace(keycloakUrl, cheURL, workspace); err != nil {
		w.Logger.Error("Failed to get user token ", zap.Error(err))
	}
	statusWorkspace, err := w.statusWorkspace(cheURL, keycloakUrl, workspace, "", testContext.WorkspaceStoppedStatus)

	if !statusWorkspace {
		return err
	}

	request, err := http.NewRequest("DELETE", cheURL+"/api/workspace/"+workspace.ID, nil)

	if err != nil {
		w.Logger.Error("Failed to delete workspace", zap.Error(err))
	}

	if keycloakAuth, err = w.KeycloakToken(keycloakUrl); err != nil {
		w.Logger.Error("Failed to get user token ", zap.Error(err))
	}

	request.Header.Add("Authorization", "Bearer "+keycloakAuth.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	_, err = w.httpClient.Do(request)

	return err
}


// StopWorkspace stop a workspace from a given workspace_id
func (w *Controller) stopWorkspace(keycloakUrl string, cheURL string, workspace *Workspace) (err error) {
	var keycloakAuth *KeycloakAuth

	request, err := http.NewRequest("DELETE", cheURL+"/api/workspace/"+workspace.ID+"/runtime", nil)

	if err != nil {
		w.Logger.Error("Failed to stop workspace", zap.Error(err))
	}

	if keycloakAuth, err = w.KeycloakToken(keycloakUrl); err != nil {
		w.Logger.Error("Failed to get user token ", zap.Error(err))
	}

	request.Header.Add("Authorization", "Bearer "+keycloakAuth.AccessToken)
	request.Header.Add("Content-Type", "application/json")

	_, err = w.httpClient.Do(request)

	return err
}
