package operator_tests

import (
	"github.com/che-incubator/che-test-harness/cmd/che/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/che-incubator/che-test-harness/cmd/che/util"
	"github.com/che-incubator/che-test-harness/docs"
	"github.com/che-incubator/che-test-harness/pkg/client"
	"github.com/che-incubator/che-test-harness/pkg/controller"
	log "github.com/che-incubator/che-test-harness/pkg/controller/logger"
	"go.uber.org/zap"

	"github.com/che-incubator/che-test-harness/pkg/monitors/metadata"
	_ "github.com/che-incubator/che-test-harness/pkg/tests"
	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	"github.com/onsi/gomega"
)

//Create Constant file
const (
	testResultsDirectory = "/test-run-results"
	jUnitOutputFilename  = "junit-che-operator.xml"
	addonMetadataName    = "addon-metadata.json"
	DebugSummaryOutput   = "debug_tests.json"
)

var Logger = &log.Zap

//SynchronizedBeforeSuite blocks are primarily meant to solve the problem of setting up the custom resources for
//Code Ready Workspaces
var _ = ginkgo.SynchronizedBeforeSuite(func() []byte {
	// Deserialize test harness configuration and assign to a struct
	if err := config.ParseConfigurationFile(); err != nil {
		Logger.Panic("Failed to get Che Test Harness Configuration. Please Check your configuration file: deploy/test-harness.yaml")
	}

	// Generate kubernetes client go to access cluster
	k8sClient, err := client.NewK8sClient()
	if err != nil {
		panic(err)
	}

	// Check if Code Ready Workspaces operator is installed on OSD namespace or external namespace
	start := util.OsdSetupNameSpace()
	if !start {
		// In case if CRW Operator not found in any namespace specified the software will crush
		os.Exit(1)
	}

	//!TODO: Try to create a specific function to call all <ginkgo suite> configuration.
	Logger.Info("Starting to setup objects before run ginkgo suite")
	// Initialize Codeready Kubernetes client to create resources in a giving namespace
	ctrl := controller.NewTestHarnessController(k8sClient)

	if !ctrl.RunTestHarness() {
		Logger.Panic("Failed to create custom resources in cluster", zap.Error(err))
	}

	return nil
}, func(data []byte) {})

var _ = ginkgo.SynchronizedAfterSuite(func() {
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
	// configure zap logging for codeready addon, Zap Logger create a file <*.log> where is possible
	//to find information about addon execution.
	Logger, _ := log.ZapLogger()

	gomega.RegisterFailHandler(ginkgo.Fail)
	Logger.Info("Code Ready Workspaces version supported: " + docs.CRW_SUPPORTED_VERSION)
	Logger.Info("Creating ginkgo reporter for Test Harness: Junit and Debug Detail reporter")

	var r []ginkgo.Reporter
	r = append(r, reporters.NewJUnitReporter(filepath.Join(testResultsDirectory, jUnitOutputFilename)))
	r = append(r, util.NewDetailsReporterFile(filepath.Join(testResultsDirectory, DebugSummaryOutput)))

	Logger.Info("Running Code Ready Workspaces e2e tests...")
	ginkgo.RunSpecsWithDefaultAndCustomReporters(t, "Code Ready Operator Test Harness", r)

	err := metadata.Instance.WriteToJSON(filepath.Join(testResultsDirectory, addonMetadataName))
	if err != nil {
		Logger.Panic("error while writing metadata")
	}
}
