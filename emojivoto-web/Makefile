svc_name := emojivoto-web

all: clean protoc webpack package-web test package

include ../common.mk

webpack:
	$(MAKE) -C webapp

package-web: webpack
	mkdir -p $(target_dir)/web
	cp web/favicon.ico $(target_dir)/web
	cp -a webapp/dist $(target_dir)

compile-vote-bot:
	CGO_ENABLED=0 go build -o $(target_dir)/emojivoto-vote-bot --trimpath cmd/vote-bot/main.go

compile: compile-vote-bot
