package idling_defaults

type Config struct {
	CheUrl               string
	KeycloakUrl          string
	IdlingTimeoutMinutes int
}

var ConfigInstance = Config{}
