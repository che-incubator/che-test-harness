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

Log into your openshift cluster, using `oc login -u <user> -p <password> <oc_api_url>.`

A properly setup Go workspace using **Go 1.13+ is required**.

Install dependencies:
```
# Install dependencies
$ go mod tidy
# Copy the dependencies to vendor folder
$ go mod vendor
# Create che-test-harness binary in bin folder. Please add the binary to the path or just execute ./bin/che-test-harness
$ make build
```

## The `che-test-harness` command

The `che-test-harness` command is the root command that executes all test harness functionality through a number of variables

### Che Test Harness Arguments

Che Test Harness comes with a number of arguments that can be passed to the `che-test-harness` command. Supported arguments:

| Argument | Usage | Default |
| -- | -- | -- |
| `--help` | Prints all available arguments | "" |
| `--che-namespace` | Namespace where Eclipe Che operator is deployed | `che` |
| `--artifacts-dir` | Directory where to store the artifacts generated by che-test-i | `/tmp/artifacts` |
| `--metrics-files` | Make reference where aws secrets are mounted  | `/etc/secrets` |

Also `che-test-harness` command support all ``Ginkgo`` flags...

# Openshift CI

Che-Test-Harness run as a part of Openshift CI every 12 hours. To visualize the jobs please go to [PROW](https://deck-ci.apps.ci.l2s4.p1.openshiftapps.com/?job=periodic-ci-che-incubator-che-test-harness-master-performance-tests).
Openshift CI Job Configuration lives in [ci-operator](https://github.com/openshift/release/tree/master/ci-operator/jobs/che-incubator/che-test-harness). How che-test-harness generate
prometheus files with test results we have to send the prom. file to s3 to use after to push the results to Prometheus PUSH Gateway.


# Workspace idling test
Test for testing if the workspace is idled after dedicated timeout.

## Setup

When launching this test, you need to know what is the time after which the workspace will be idled. By default, workspace idle timeout is set to 30 minutes. You can change that by adding following to your `CheCluster`:

```
spec:
  server:
    customCheProperties:
      CHE_LIMITS_WORKSPACE_IDLE_TIMEOUT: "<your timeout in miliseconds>"
```

A properly setup Go workspace using **Go 1.13+ is required**.

Install dependencies:
```
# Install dependencies
$ go mod tidy
# Copy the dependencies to vendor folder
$ go mod vendor
# Create che-test-idling binary in bin folder. Please add the binary to the path or just execute ./bin/che-test-idling
$ make build-idling
```

## The `che-test-idling` command

The `che-test-idling` command is the root command that executes all required functionality fort testing the idling.

### Che Test Idling Arguments

Che Test Idling comes with a number of arguments that can be passed to the `che-test-idling` command. Supported arguments:

| Argument | Usage | Default |
| -- | -- | -- |
| `--help` | Prints all available arguments | "" |
| `--username` | Username that test will use for login into the Che | `admin` |
| `--password` | Password that test will use for login into the Che | `admin` |
| `--che-url` | URL of running Che instance | | 
| `--idling-timeout` | Timeout in MINUTES after which workspace should be stopped. | 30 |

Also `che-test-idling` command support all ``Ginkgo`` flags.

**Note:** After the idling timeout, workspace is "stopped" - but it take some extra time until the workspace really becomes stopped. That means that the workspace is stopped a little bit later then the idling timeout says. Average extra time is about 2-3 minutes, so the test is adding extra 5 minutes for the workspace to be idled (just to be on a safe side).