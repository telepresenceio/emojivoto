.PHONY: package protoc test

target_dir := target

clean:
	rm -rf gen
	rm -rf $(target_dir)
	mkdir -p $(target_dir)
	mkdir -p gen

protoc:
	protoc -I .. ../proto/*.proto --go_out=paths=source_relative:./gen --go-grpc_out=paths=source_relative:./gen

package: protoc compile build-container

build-container:
	docker build .. -t "$(IMAGE_REGISTRY)/$(svc_name):$(IMAGE_TAG)" --build-arg svc_name=$(svc_name)

build-multi-arch:
	docker buildx build .. -t "$(IMAGE_REGISTRY)/$(svc_name):$(IMAGE_TAG)" --build-arg svc_name=$(svc_name) \
		-f ../Dockerfile-multi-arch --platform linux/amd64,linux/arm64,linux/arm/v7 --push

compile:
	CGO_ENABLED=0 go build -o $(target_dir)/$(svc_name) --trimpath cmd/server.go

test:
	go test ./...

run:
	go run cmd/server.go
