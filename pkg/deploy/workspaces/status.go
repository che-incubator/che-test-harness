package workspaces

import (
	"encoding/json"
	"errors"
	"github.com/che-incubator/che-test-harness/pkg/client"
	testContext "github.com/che-incubator/che-test-harness/pkg/deploy/context"
	"github.com/che-incubator/che-test-harness/pkg/monitors"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (w *Controller) statusWorkspace(cheURL string, keycloakURL string, workspace *Workspace, workspaceStackName string, desiredStatus string) (boolean bool, err error) {
	var keycloakAuth *KeycloakAuth
	timeout := time.After(5 * time.Minute)
	tick := time.Tick(15 * time.Second)

	if desiredStatus == testContext.WorkspaceRunningStatus {
		stopCh := make(chan struct{})
		k8sClient, err := client.NewK8sClient()
		if err != nil {
			w.Logger.Panic("Failed to create kubernetes client.", zap.Error(err))
		}
		defer close(stopCh)
		monitor, _ := monitors.NewMonitor(k8sClient.Kube())
		go func() {
			if err := monitor.DescribeEvents(stopCh, workspaceStackName); err != nil {
				panic(err)
			}
		}()
	}

	for {
		select {
		case <-timeout:
			return false, errors.New("workspace didn't start after 5 minutes")
		case <-tick:
			if keycloakAuth, err = w.KeycloakToken(keycloakURL); err != nil {
				w.Logger.Error("Failed to get user token ", zap.Error(err))
			}
			request, err := http.NewRequest("GET", cheURL+"/api/workspace/"+workspace.ID, nil)
			if err != nil {
				return false, err
			}
			request.Header.Add("Authorization", "Bearer "+keycloakAuth.AccessToken)
			request.Header.Add("Content-Type", "application/json")

			response, err := w.httpClient.Do(request)
			err = json.NewDecoder(response.Body).Decode(&workspace)

			if workspace.Status == 	desiredStatus {
				return true, nil
			}
		}
	}
}
