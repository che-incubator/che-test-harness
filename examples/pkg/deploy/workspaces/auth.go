package workspaces

import (
	"encoding/json"
	"github.com/che-incubator/che-test-harness/pkg/deploy/context"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strings"
)

// KeycloakToken return a JWT from keycloak
func (w *Controller) KeycloakToken(keycloakTokenUrl string) (keycloakAuth *KeycloakAuth, err error) {
	data := url.Values{
		"client_id"  : {ClientID},
		"username"   : {context.TestInstance.Setup.Username},
		"password"   : {context.TestInstance.Setup.Password},
		"grant_type" : {"password"},
	}
	request, err := http.NewRequest("POST", keycloakTokenUrl, strings.NewReader(data.Encode()))

	if err != nil {
		w.Logger.Panic("Failed to get token from keycloak", zap.Error(err))
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response, err := w.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&keycloakAuth)

	return keycloakAuth, err
}
