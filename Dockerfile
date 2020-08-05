FROM registry.svc.ci.openshift.org/openshift/release:golang-1.13 AS builder

ENV PKG=/go/src/github.com/quay.io/operator-tests/
WORKDIR ${PKG}

# compile test binary
COPY . .
RUN make

FROM registry.access.redhat.com/ubi7/ubi-minimal:latest

ENV CODEREADY_NAMESPACE=codeready-workspaces-operator-qe

COPY --from=builder /go/src/github.com/quay.io/operator-tests/bin/che-test-harness che-test-harness

ENTRYPOINT [ "/che-test-harness" ]
