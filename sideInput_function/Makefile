.PHONY: build
build:clean
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/sideinput-example main.go

.PHONY: image
image: build
	docker buildx build --no-cache -t "quay.io/numaio/sideinput-sink-e2e-test:latest" --platform linux/amd64,linux/arm64 --target sideinput . --push

clean:
	-rm -rf ./dist