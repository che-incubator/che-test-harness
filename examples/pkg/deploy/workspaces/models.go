package workspaces

type KeycloakAuth struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Workspace struct {
	ID         string     `json:"id"`
	Attributes Attributes `json:"attributes"`
	Status     string     `json:"status"`
}

type Attributes struct {
	InfrastructureNamespace string `json:"infrastructureNamespace"`
}

