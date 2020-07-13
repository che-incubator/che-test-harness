# Run Code Ready Workspaces Test Harness in OSD
To run test harness in OSD first we need to have access to OSD clusters and access to Code Ready Workspaces addon. The
addons are managed in `managed-tenants repo`.

1. Get OFFLINE_TOKEN from OSD (EG: cloud.redhat.com/openshift/token) and put it into `execute-test-harness.sh` file.
2. Get cluster ID from OSD and put it into `osd-test-harness.sh` file
3. Run launch script.

    ```
    ./osd-test-harness.sh
    ```

# Run Code Ready Workspaces Test Harness outside of OSD
1. Access To Openshift Cluster

2. Login to the cluster as `kube:admin`

   ```
   oc login ...
   ```

3. Run the test from your machine

   ```
   ./run-crw-testsuite.sh <namespace> <report-dir> [!<crw-operator-namespace>]
   ```

Where are:
 - `namespace` - namespace where you want to deploy test-harness.
 - `report-dir` - Directory where you want to download the results of tests from pods.
 - `crw-operator-namespace` - Namespace where Code Ready Workspaces operator is deployed. 
If you will deploy test-harness in the same namespace with Code Ready Workspace Operator this option is not required.
