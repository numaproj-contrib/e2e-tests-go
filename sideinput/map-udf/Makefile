.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/map-sideinput-udf main.go

.PHONY: image
#To build an image that supports multiple platforms(linux/amd64,linux/arm64) and push to quay.io, use the following command
image: build
	docker buildx build -t "quay.io/numaio/numaproj-contrib/e2e-map-sideinput-udf:v0.0.10" --platform linux/amd64,linux/arm64 --target map-sideinput-udf . --push

clean:
	-rm -rf ./dist