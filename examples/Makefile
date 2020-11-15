OUT_FILE := ./bin/che-performance-test

build-performance:
	CGO_ENABLED=0 go test -v -c -o ${OUT_FILE} ./cmd/che-performance/che_performance_test.go
