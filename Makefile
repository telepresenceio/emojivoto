IMAGE_REGISTRY ?= ghcr.io/telepresenceio
IMAGE_TAG ?= 0.2.0
export IMAGE_REGISTRY
export IMAGE_TAG

.PHONY: web emoji voting test push kustomize

all: build test

web:
	$(MAKE) -C emojivoto-web

emoji:
	$(MAKE) -C emojivoto-emoji

voting:
	$(MAKE) -C emojivoto-voting

build: web emoji voting

clean:
	$(MAKE) -C emojivoto-web clean
	$(MAKE) -C emojivoto-emoji clean
	$(MAKE) -C emojivoto-voting clean

push:
	$(MAKE) -C emojivoto-web build-multi-arch
	$(MAKE) -C emojivoto-emoji build-multi-arch
	$(MAKE) -C emojivoto-voting build-multi-arch

deploy-to-cluster: kustomize/deployment/kustomization.yaml
	@kubectl apply -k $(<D)

deploy-to-docker-compose: compose.yaml
	docker compose stop
	docker compose rm -vf
	$(MAKE) -C emojivoto-web build-container
	$(MAKE) -C emojivoto-emoji build-container
	$(MAKE) -C emojivoto-voting build-container
	docker compose up -d

compose.yaml: compose.yaml.in FORCE
	@envsubst < $< > $@

kustomize/deployment/kustomization.yaml: kustomize/deployment/kustomization.yaml.in FORCE
	@envsubst < $< > $@
FORCE:

kustomize: kustomize/deployment/kustomization.yaml
	@kubectl kustomize $(<D)

test:
	go test ./...
