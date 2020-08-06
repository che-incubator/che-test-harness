package operator_tests

import (
	"flag"
	"fmt"
	"github.com/che-incubator/che-test-harness/pkg/common"
	"github.com/che-incubator/che-test-harness/pkg/common/aws"
	"github.com/che-incubator/che-test-harness/pkg/common/client"
	"github.com/che-incubator/che-test-harness/pkg/common/logger"
	"github.com/che-incubator/che-test-harness/pkg/common/prometheus"
	"github.com/che-incubator/che-test-harness/pkg/common/reporter"
	"github.com/che-incubator/che-test-harness/pkg/controller"
	"github.com/che-incubator/che-test-harness/pkg/monitors/metadata"
	_ "github.com/che-incubator/che-test-harness/pkg/suites/che"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"path/filepath"
	"testing"
)

//Create Constant file
const (
	jUnitOutputFilename = "junit-che-operator.xml"
	addonMetadataName   = "addon-metadata.json"
	DebugSummaryOutput  = "debug_tests.json"
)

// Create an instance to save all our data there
type Configs struct {
	artifactsDir string
	metricsFiles string
}

var CfgInstance = Configs{}

// Start to register flags
func init()  {
	registerCheFlags(flag.CommandLine)
}

// SynchronizedBeforeSuite blocks are primarily meant to solve the problem of setting up the custom resources for Eclipse Che
var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	// Generate kubernetes client go to access cluster
	var Logger, err = logger.ZapLogger()
	if err != nil {
		panic("Failed to create zap logger")
	}
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		panic("Failed to create k8s client go")
	}

	// Initialize Kubernetes client to create resources in a giving namespace
	ctrl := controller.NewTestHarnessController(k8sClient)

	if !ctrl.RunTestHarness() {
		Logger.Panic("Failed to create custom resources in cluster", zap.Error(err))
	}

	return nil
}, func(data []byte) {})

var _ = ginkgo.SynchronizedAfterSuite(func() {
	var Logger, _ = logger.ZapLogger()

	_ = metadata.Instance.WriteToJSON(filepath.Join(CfgInstance.artifactsDir, addonMetadataName))

	newMetrics := prometheus.NewMetrics()
	if newMetrics == nil {
		Logger.Panic("Error getting new prometheus provider")
	}

	// Generate a new prometheus files from
	prometheusFilename, err := newMetrics.WritePrometheusFile(CfgInstance.artifactsDir)
	if err != nil {
		Logger.Panic("Error while writing prometheus prometheus", zap.Error(err))
	}

	if len(CfgInstance.metricsFiles) > 0 {
		if err := sendDataToS3(prometheusFilename); err != nil {
			Logger.Panic("Error sending data to aws s3", zap.Error(err))
		}
	}

	k8sClient, err := client.NewK8sClient()
	if err != nil {
		panic(err)
	}

	ctrl := controller.NewTestHarnessController(k8sClient)

	//Delete all objects after pass all test suites.
	Logger.Info("Clean up all created objects by Test Harness.")

	if err := ctrl.DeleteCustomResource(); err != nil {
		Logger.Panic("Failed to remove Kubernetes Puller Image from Cluster")
	}

	if err := ctrl.DeleteNamespace(); err != nil {
		Logger.Panic("Failed to remove Kubernetes Puller Image from Cluster")
	}
}, func() {})

func TestHarnessChe(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)

	var r []ginkgo.Reporter
	r = append(r, reporters.NewJUnitReporter(filepath.Join(CfgInstance.artifactsDir, jUnitOutputFilename)))
	r = append(r, reporter.NewDetailsReporterFile(filepath.Join(CfgInstance.artifactsDir, DebugSummaryOutput)))

	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "Eclipse Che Test Harness", r)

}

// Register All flags used by test harness
func registerCheFlags(flags *flag.FlagSet) {
	flags.StringVar(&CfgInstance.artifactsDir, "artifacts-dir", "/tmp/artifacts", "If is specified test harness will save all reports in the given directory, if not will save artifacts in the current directory. Default dir is /tmp/artifacts")
	flags.StringVar(&CfgInstance.metricsFiles, "metrics-files", "", "If it is set che test harness start to send the data to AWS S3 . You should have valid secrets in the files")
	flags.StringVar(&metadata.Namespace.Name, "che-namespace", "eclipse-che", "Namespace where che-operator was deployed before launch tests. Default namespace is `eclipse-che`")
}

// Generate a Prometheus file from json metadata and connect with aws s3 and put the data into buckets
func sendDataToS3(prometheusFilename string) error {
	if err := common.LoadConfigs(CfgInstance.metricsFiles); err != nil {
		return fmt.Errorf("error loading initial state: %v", err)
	}

	if err := uploadFileToMetricsBucket(filepath.Join(CfgInstance.artifactsDir, prometheusFilename)); err != nil {
		return fmt.Errorf("error while uploading prometheus prometheus: %v", err)
	}

	return nil
}

// uploadFileToMetricsBucket uploads the given file (with absolute path) to the prometheus S3 bucket "incoming" directory.
func uploadFileToMetricsBucket(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	return aws.WriteToS3(aws.CreateS3URL(viper.GetString("prometheus-datahub"), "prometheus-datahub", filepath.Base(filename)), data)
}
