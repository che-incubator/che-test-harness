package idling_test

import (
	"flag"
	"strings"
	"testing"

	"github.com/onsi/ginkgo"

	idling_defaults "github.com/che-incubator/che-test-harness/pkg/common/idling"
	"github.com/che-incubator/che-test-harness/pkg/monitors/metadata"
	_ "github.com/che-incubator/che-test-harness/pkg/suites/workspace-idling"
	"github.com/onsi/gomega"
)

// Start to register flags
func init() {
	registerCheFlags(flag.CommandLine)
}

func setupKeycloakUrl() {
	idling_defaults.ConfigInstance.KeycloakUrl = strings.Replace(idling_defaults.ConfigInstance.CheUrl, "che", "keycloak", 1) + "/auth/realms/che/protocol/openid-connect/token"
}

// Register All flags used by idling test
func registerCheFlags(flags *flag.FlagSet) {
	flags.StringVar(&metadata.User.Username, "username", "admin", "Username of a testing user that can log in Che and create workspaces. Defaults to 'admin'.")
	flags.StringVar(&metadata.User.Password, "password", "admin", "Password of a testing user that can log in Che and create workspaces. Defaults to 'admin'.")
	flags.StringVar(&idling_defaults.ConfigInstance.CheUrl, "che-url", "eclipse-che-url", "URL of Che deployment.")
	flags.IntVar(&idling_defaults.ConfigInstance.IdlingTimeoutMinutes, "idling-timeout", 30, "Timeout after which workspace should be idled. Default to 30 minutes.")
}

func TestWorkspaceIdling(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	setupKeycloakUrl()

	ginkgo.RunSpecs(t, "Workspace Idling")
}
