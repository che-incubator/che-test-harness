package context

type Setup struct {
	ArtifactsDir   string `json:"artifacts_dir"`
	AwsSecretFiles string `json:"aws_secret_folder,omitempty" binding:"required"`
	CheNamespace   string `json:"che_namespace"`
	DeployChe      bool   `json:"deploy_che"`
	Username       string `json:"username"`
	Password       string `json:"password"`
}

