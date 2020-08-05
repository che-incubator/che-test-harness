package aws

import (
	"fmt"
	"github.com/che-incubator/che-test-harness/pkg/common"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/prometheus/common/log"
	"github.com/spf13/viper"
)

const (
	metricsAWSAccessKeyID     = "metrics.awsAccessKeyId"
	metricsAWSSecretAccessKey = "metrics.awsSecretAccessKey"
	metricsAWSRegion          = "metrics.awsRegion"

	metricsAWSAccessKeyIDEnv     = "METRICS_AWS_ACCESS_KEY_ID"
	metricsAWSSecretAccessKeyEnv = "METRICS_AWS_SECRET_ACCESS_KEY"
	metricsAWSRegionEnv          = "METRICS_AWS_REGION"
)

type awsSession struct {
	session *session.Session
	once    sync.Once
}

// AWSSession is the global AWS session for interacting with S3.
var AWSSession awsSession

func init() {
	viper.BindEnv(metricsAWSAccessKeyID, metricsAWSAccessKeyIDEnv)
	common.RegisterSecret(metricsAWSAccessKeyID, "metrics-aws-access-key")

	viper.BindEnv(metricsAWSSecretAccessKey, metricsAWSSecretAccessKeyEnv)
	common.RegisterSecret(metricsAWSSecretAccessKey, "metrics-aws-secret-access-key")

	viper.BindEnv(metricsAWSRegion, metricsAWSRegionEnv)
	common.RegisterSecret(metricsAWSRegion, "metrics-aws-region")
}

// Parse creds from viper and create aws session
func (a *awsSession) getSession() (*session.Session, error) {
	var err error
	a.once.Do(func() {
		a.session, err = session.NewSession(aws.NewConfig().
			WithCredentials(credentials.NewStaticCredentials(viper.GetString(metricsAWSAccessKeyID), viper.GetString(metricsAWSSecretAccessKey), "")).
			WithRegion(viper.GetString(metricsAWSRegion)))

		if err != nil {
			log.Errorf("error initializing AWS session: %v", err)
		}
	})

	if a.session == nil {
		err = fmt.Errorf("unable to initialize AWS session")
	}

	return a.session, err
}
