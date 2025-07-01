include .env
export

.PHONY: web emoji-svc voting-svc test push kustomize

all: build test

web:
	$(MAKE) -C emojivoto-web

emoji-svc:
	$(MAKE) -C emojivoto-emoji-svc

voting-svc:
	$(MAKE) -C emojivoto-voting-svc

build: web emoji-svc voting-svc

multi-arch:
	$(MAKE) -C emojivoto-web build-multi-arch
	$(MAKE) -C emojivoto-emoji-svc build-multi-arch
	$(MAKE) -C emojivoto-voting-svc build-multi-arch

deploy-to-minikube:
	$(MAKE) -C emojivoto-web build-container
	$(MAKE) -C emojivoto-emoji-svc build-container
	$(MAKE) -C emojivoto-voting-svc build-container
	kubectl delete -f emojivoto.yml || echo "ok"
	kubectl apply -f emojivoto.yml

deploy-to-docker-compose:
	docker compose stop
	docker compose rm -vf
	$(MAKE) -C emojivoto-web build-container
	$(MAKE) -C emojivoto-emoji-svc build-container
	$(MAKE) -C emojivoto-voting-svc build-container
	docker compose up -d

push-%:
	docker push $(IMAGE_REGISTRY)/emojivoto-$*:$(IMAGE_TAG)

push: multi-arch push-emoji-svc push-voting-svc push-web

kustomize/deployment/kustomization.yml: kustomize/deployment/kustomization.yml.in .env
	@envsubst < $< > $@

kustomize: kustomize/deployment/kustomization.yml
	@kubectl kustomize $(<D)

test:
	go test ./...
