DOCKERIO_ORG=quay.io/numaio
PLATFORM=linux/x86_64
TARGET=redis

IMAGE_TAG=$(TAG)
ifeq ($(IMAGE_TAG),)
	IMAGE_TAG="latest"
endif

.PHONY: build
build: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/simple-source-example main.go

.PHONY: image

image: build
	docker buildx build --no-cache -t "quay.io/numaio/numaflow-source/simple_source_with_sideinput:latest"  --platform $(PLATFORM)  --target simple-source . --push

.PHONY: lint
lint:
	go mod tidy
	golangci-lint run --fix --verbose --concurrency 4 --timeout 5m

clean:
	-rm -rf ./dist
