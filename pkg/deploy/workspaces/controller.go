package workspaces

import (
	"github.com/che-incubator/che-test-harness/internal/logger"
	"net/http"
)

// WorkspacesController useful to add logger and http client.
type Controller struct {
	httpClient *http.Client
	Logger     logger.Log
}

// NewWorkspaceController creates a new WorkspacesController from a given client.
func NewWorkspaceController(c *http.Client) *Controller {
	return &Controller{
		httpClient: c,
		Logger:     logger.Zap,
	}
}
