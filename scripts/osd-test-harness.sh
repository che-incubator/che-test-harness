#!/usr/bin/env bash
set -e

export OCM_TOKEN="<OFFLINE_TOKEN>"
export ADDON_IDS="codeready-workspaces addon name"
export REPORT_DIR="Report_DIR to throw results"
export CLUSTER_ID="CLUSTER_ID where you want to execute the tests in OSD"
export ADDON_TEST_HARNESSES=quay.io/crw/osd-e2e:latest # change to nightly to check crw dev product

osde2e test -configs stage,addon-suite,skip-health-checks
