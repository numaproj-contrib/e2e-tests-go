.PHONY: build
build:clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/redis-e2e-test-sink main.go

.PHONY: image
image: build
	docker buildx build --no-cache -t "quay.io/numaio/numaflow-sink/redis-e2e-test-sink-with-sideinput:latest" --platform linux/amd64,linux/arm64 --target redis . --push


clean:
	-rm -rf ./dist