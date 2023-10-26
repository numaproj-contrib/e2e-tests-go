DOCKERIO_ORG=quay.io/numaio
PLATFORM=linux/x86_64,linux/amd64,linux/arm64
TARGET=redis

IMAGE_TAG=$(TAG)
ifeq ($(IMAGE_TAG),)
	IMAGE_TAG=latest
endif

.PHONY: build
build:clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/sideinput-example main.go

.PHONY: image
image: build
	docker buildx build --no-cache -t "quay.io/numaio/numaflow-sideinput/sideinput-sink-e2e-test:latest" --platform $(PLATFORM) --target sideinput . --push

.PHONY: lint
lint:
	go mod tidy
	golangci-lint run --fix --verbose --concurrency 4 --timeout 5m

clean:
	-rm -rf ./dist