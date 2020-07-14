OUT_FILE := ./bin/che-test-harness
DOCKER_IMAGE_NAME :=quay.io/crw/osd-e2e
CODE_READY_VERSION:=$(shell grep 'CODE_READY_VERSION' docs/version.go | awk '{ print $$4 }' | tr -d '"')
CODE_READY_NIGHTLY:=$(shell grep 'CODE_READY_NIGHTLY' docs/version.go | awk '{ print $$4 }' | tr -d '"')
IS_NIGHTLY:=$(shell grep 'IS_NIGHTLY' docs/version.go | awk '{ print $$4 }' | tr -d '"')

ifeq ($(IS_NIGHTLY),true)
	CODE_READY_VERSION=nightly
else
	CODE_READY_VERSION=latest
endif

build:
	CGO_ENABLED=0 go test -v -c -o ${OUT_FILE} ./cmd/che/che_harness_test.go

build-container:
	docker build -t $(DOCKER_IMAGE_NAME):$(CODE_READY_VERSION) --no-cache .

push-container:
	docker push $(DOCKER_IMAGE_NAME):$(CODE_READY_VERSION)
