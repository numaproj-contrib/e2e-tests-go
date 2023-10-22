.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./dist/redis-e2e-test-sink main.go

.PHONY: image
image: build
	docker buildx build --no-cache -t "quay.io/numaio/redis-e2e-test-sink-sideinput:latest" --target redis .

clean:
	-rm -rf ./dist