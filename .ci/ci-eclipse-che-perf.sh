#!/bin/bash

export CODEREADY_NAMESPACE="che"

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
    oc create namespace ${CODEREADY_NAMESPACE}
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/service_account.yaml -n ${CODEREADY_NAMESPACE}
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/crds/org_v1_che_crd.yaml
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/role.yaml -n ${CODEREADY_NAMESPACE}
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/role_binding.yaml -n ${CODEREADY_NAMESPACE}
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/cluster_role.yaml
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/cluster_role_binding.yaml
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/role_binding_oauth.yaml
    oc apply -f https://raw.githubusercontent.com/eclipse/che-operator/master/deploy/operator.yaml -n ${CODEREADY_NAMESPACE}
}

function deployTestHArness() {
    make build
    ${TEST_HARNESS_ROOT}/bin/che-test-harness
}

function run() {
    init
    installCheOperator
    deployTestHArness
}

run
