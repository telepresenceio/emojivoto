IMAGE_REGISTRY ?= ghcr.io/telepresenceio
IMAGE_TAG ?= 0.1.1
export IMAGE_REGISTRY
export IMAGE_TAG

.PHONY: web emoji-svc voting-svc test push kustomize

all: build test

web:
	$(MAKE) -C emojivoto-web

emoji-svc:
	$(MAKE) -C emojivoto-emoji-svc

voting-svc:
	$(MAKE) -C emojivoto-voting-svc

build: web emoji-svc voting-svc

push:
	$(MAKE) -C emojivoto-web build-multi-arch
	$(MAKE) -C emojivoto-emoji-svc build-multi-arch
	$(MAKE) -C emojivoto-voting-svc build-multi-arch

deploy-to-minikube:
	$(MAKE) -C emojivoto-web build-container
	$(MAKE) -C emojivoto-emoji-svc build-container
	$(MAKE) -C emojivoto-voting-svc build-container
	kubectl delete -f emojivoto.yaml || echo "ok"
	kubectl apply -f emojivoto.yaml

deploy-to-docker-compose: compose.yaml
	docker compose stop
	docker compose rm -vf
	$(MAKE) -C emojivoto-web build-container
	$(MAKE) -C emojivoto-emoji-svc build-container
	$(MAKE) -C emojivoto-voting-svc build-container
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
