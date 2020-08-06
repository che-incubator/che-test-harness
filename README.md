# Che-Test-Harness
Testing solution written in golang using ginkgo framework for Eclipse Che. This tests runs in Openshift CI Platform. 

# Specifications
* Instrumented tests with ginkgo framework. Find more info: https://onsi.github.io/ginkgo/
* Structured logging with zap.
* Use client-go to connect to Openshift Cluster.
* Deploy Eclipse Che nightly in Cluster.
* Defined events watcher oriented to Eclipse Che Resources. Please look `pkg/monitors/watcher.go`
* Deploy Kubernetes Image Puller in Cluster which will pre-pull workspaces images.
* Create, start and get measure up times of Eclipse Che Workspaces
* Transform Json results of tests in prometheus language and send this file to AWS S3 to be consumed for Prometheus Push Gateway if aws will be provided

# Setup

This is an example test harness meant for testing the che operator addon. It does the following:

* Tests for the existence of CRD in cluster. This should be present if the che
  operator addon has been installed properly.
 * Check the pods health
* Writes out a junit XML file with tests results to the /test-run-results directory as expected
  by the [https://github.com/openshift/osde2e](osde2e) test framework.
* Writes out an `addon-metadata.json` file which will also be consumed by the osde2e test framework.
# Tests execution

In order to tests locally osde2e tests you should execute first `make build` which create a new
binary in ./bin folder.
