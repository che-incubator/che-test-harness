#!/bin/bash

set -ex

export METRICS_FILES="/usr/local/ci-secrets/test-harness-secrets"
export ARTIFACTS_DIR="/tmp/artifacts"
export OPERATOR_NAMESPACE="che"

function init() {
  # shellcheck disable=SC2155
  local SCRIPT=$(readlink -f "$0")
  # shellcheck disable=SC2155
  local SCRIPT_DIR=$(dirname "$SCRIPT")
  if [[ ${WORKSPACE} ]] && [[ -d ${WORKSPACE} ]]; then
    export TEST_HARNESS_ROOT=${WORKSPACE};
  else
    # shellcheck disable=SC2155
    export TEST_HARNESS_ROOT=$(dirname "$SCRIPT_DIR");
  fi
}

function installCheOperator() {
    oc create namespace ${OPERATOR_NAMESPACE}
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/service_account.yaml -n ${OPERATOR_NAMESPACE}
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/crds/org_v1_che_crd.yaml
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/role.yaml -n ${OPERATOR_NAMESPACE}
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/role_binding.yaml -n ${OPERATOR_NAMESPACE}
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/cluster_role.yaml
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/cluster_role_binding.yaml
    #oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/role_binding_oauth.yaml
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/operator.yaml -n ${OPERATOR_NAMESPACE}
}

function deployTestHArness() {
    # For some reason go on PROW force usage vendor folder
    # This workaround is here until we don't figure out cause
    go mod tidy
    go mod vendor

    make  build-performance
    "${TEST_HARNESS_ROOT}"/bin/che-performance-test --che-namespace=${OPERATOR_NAMESPACE} --metrics-files=${METRICS_FILES} --artifacts-dir=${ARTIFACTS_DIR}

}

function run() {
    init
    installCheOperator
    deployTestHArness
}

echo "${BUILD_ID}"
echo "${BUILD_NUMBER}"

run
