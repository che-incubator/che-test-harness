OUT_FILE := ./bin/che-test-harness
IDLE_OUT_FILE := ./bin/che-test-idling

build:
	CGO_ENABLED=0 go test -v -c -o ${OUT_FILE} ./cmd/che/che_harness_test.go

build-idling:
	CGO_ENABLED=0 go test -v -c -o ${IDLE_OUT_FILE} ./cmd/idling/workspace_idling_test.go
