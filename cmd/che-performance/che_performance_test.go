package performance_test

import (
	"flag"
	"path/filepath"
	"testing"

	"github.com/che-incubator/che-test-harness/internal/logger"
	reporter "github.com/che-incubator/che-test-harness/internal/reporters"
	"github.com/che-incubator/che-test-harness/pkg/client"
	"github.com/che-incubator/che-test-harness/pkg/deploy"
	"github.com/onsi/ginkgo/reporters"
	"go.uber.org/zap"

	"github.com/che-incubator/che-test-harness/pkg/deploy/context"
	_ "github.com/che-incubator/che-test-harness/test/performance"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	jUnitOutputFilename = "junit-performance-tests.xml"
	instanceMetadata    = "che-performance.json"
	DebugSummaryOutput  = "debug_tests.json"
)

// Start to register flags
func init() {
	registerCheFlags(flag.CommandLine)
}

// Register All flags used by test harness
func registerCheFlags(flags *flag.FlagSet) {
	flags.StringVar(&context.TestInstance.ArtifactsDir, "artifacts-dir", "/home/flacatus/WORKSPACE/incubator/che-test-harness/build", "If is specified test harness will save all reports in the given directory, if not will save artifacts in the current directory. Default dir is /tmp/artifacts")
	flags.StringVar(&context.TestInstance.AwsSecretFiles, "aws-secrets-folder", "/tmp", "If it is set che test harness start to send the data to AWS S3. You should have valid aws secrets in the folder")
	flags.StringVar(&context.TestInstance.Setup.CheNamespace, "che-namespace", "che", "Namespace where che-operator was deployed before launch tests. Default namespace is `eclipse-che`")
	flags.BoolVar(&context.TestInstance.Setup.DeployChe, "deploy-che", true, "Deploy Eclipse Che into cluster. By default is true")
	flags.StringVar(&context.TestInstance.Setup.Username, "username", "admin", "Username of user to log in. Default username is `admin`")
	flags.StringVar(&context.TestInstance.Setup.Password, "password", "admin", "Password of user to log in. Default password is `admin`")
}

var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	var Logger, err = logger.ZapLogger()
	if err != nil {
		panic("Failed to create zap logger")
	}

	k8sClient, err := client.NewK8sClient()
	if err != nil {
		Logger.Panic("Failed to create k8s client go", zap.Error(err))
	}

	if context.TestInstance.DeployChe {
		// Initialize Kubernetes client to create resources in a giving namespace
		ctrl := deploy.NewTestHarnessController(k8sClient)

		if !ctrl.RunTestHarness() {
			Logger.Panic("Failed to deploy Eclipse Che", zap.Error(err))
		}
	}

	return nil
}, func(data []byte) {})

var _ = ginkgo.SynchronizedAfterSuite(func() {
	_ = context.TestInstance.WriteToJSON(filepath.Join(context.TestInstance.ArtifactsDir, instanceMetadata))

	var Logger, _ = logger.ZapLogger()

	k8sClient, err := client.NewK8sClient()
	if err != nil {
		panic(err)
	}

	ctrl := deploy.NewTestHarnessController(k8sClient)
	//Delete all objects after pass all test suites.
	Logger.Info("Clean up all created objects by Test Harness.")

	if err := ctrl.DeleteCustomResource(); err != nil {
		Logger.Panic("Failed to remove Kubernetes Puller Image from Cluster")
	}
}, func() {})

func TestChePerformance(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)

	var r []ginkgo.Reporter
	r = append(r, reporters.NewJUnitReporter(filepath.Join(context.TestInstance.ArtifactsDir, jUnitOutputFilename)))
	r = append(r, reporter.NewDetailsReporterFile(filepath.Join(context.TestInstance.ArtifactsDir, DebugSummaryOutput)))
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "Eclipse Che Performance Tests", r)

}
