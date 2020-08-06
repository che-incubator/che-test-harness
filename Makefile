OUT_FILE := ./bin/che-test-harness

build:
	CGO_ENABLED=0 go test -v -c -o ${OUT_FILE} ./cmd/che/che_harness_test.go
