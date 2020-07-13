# che-operator-test-harness

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
