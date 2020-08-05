#!/usr/bin/env bash

set -e

CRW_TEST_NAMESPACE=$1
REPORT_DIR=$2
CRW_OPERATOR_NAMESPACE=$3

if [[ "${CRW_TEST_NAMESPACE}" == "" || ${REPORT_DIR} == ""  ]]; then
  echo "Please specify namespace as the first argument to run test harness."
  echo "crw-operator-namespace argument is not necessary only if codeready workspaces is deployed"
  echo "into a distint namespace where will deploy our test harness"
  echo "execute-test-harness.sh <namespace> <report-dir> [!<crw-operator-namespace>]"
  exit 1
fi

if [ "${CRW_OPERATOR_NAMESPACE}" == "" ]; then
  CRW_OPERATOR_NAMESPACE=CRW_TEST_NAMESPACE
fi

ID=$(date +%s)
OPENSHIFT_API_URL=$(oc secrets view --minify -o jsonpath='{.clusters[*].cluster.server}')
OPENSHIFT_API_TOKEN=$(oc whoami -t)

TMP_POD_YML=$(mktemp)
TMP_KUBECONFIG_YML=$(mktemp)

cat kubeconfig.template.yml |
    sed -e "s#__OPENSHIFT_API_URL__#${OPENSHIFT_API_URL}#g" |
    sed -e "s#__OPENSHIFT_API_TOKEN__#${OPENSHIFT_API_TOKEN}#g" |
    cat >${TMP_KUBECONFIG_YML}

cat ${TMP_KUBECONFIG_YML}

oc delete configmap -n ${CRW_TEST_NAMESPACE} crw-testsuite-kubeconfig || true
echo "A"
oc create configmap -n ${CRW_TEST_NAMESPACE} crw-testsuite-kubeconfig \
    --from-file=secrets=${TMP_KUBECONFIG_YML}

cat test-harness.pod.template.yml |
    sed -e "s#__ID__#${ID}#g" |
    sed -e "s#__NAMESPACE__#${CRW_TEST_NAMESPACE}#g" |
    sed -e "s#__CODEREADY_NAMESPACE__#${CRW_TEST_NAMESPACE}#g" |
    cat >${TMP_POD_YML}

cat ${TMP_POD_YML}

# start the test
oc create -f ${TMP_POD_YML}

# wait for the pod to start
while true; do
    sleep 3
    PHASE=$(oc get pod -n ${CRW_TEST_NAMESPACE} crw-testsuite-${ID} \
        --template='{{ .status.phase }}')
    if [[ ${PHASE} == "Running" ]]; then
        break
    fi
done

# wait for the test to finish
oc logs -n ${CRW_TEST_NAMESPACE} crw-testsuite-${ID} -c test -f

# just to sleep
sleep 3

# download the test results
mkdir -p ${REPORT_DIR}/${ID}

oc rsync -n ${CRW_TEST_NAMESPACE} \
    crw-testsuite-${ID}:/test-run-results ${REPORT_DIR}/${ID} -c download

oc exec -n ${CRW_TEST_NAMESPACE} crw-testsuite-${ID} -c download \
    -- touch /tmp/done
