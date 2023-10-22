.PHONY: build
build: clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/simple-source-example main.go

.PHONY: image
image: build
	docker buildx build --no-cache -t "quay.io/numaio/simple_source_with_sideinput:latest"  --platform linux/amd64,linux/arm64  --target simple-source . --push

clean:
	-rm -rf ./dist
