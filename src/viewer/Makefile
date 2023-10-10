PACKAGE_ROOT := github.com/tsubasaogawa/lambda-image-viewer/src/viewer

.PHONY: compile compile-viewer compile-tagger deploy default
.DEFAULT_GOAL := default

compile:
	make compile-viewer compile-tagger

compile-viewer: cmd/viewer/*.go
	GOARCH=amd64 GOOS=linux go build -o bin/main $(PACKAGE_ROOT)/cmd/viewer

compile-tagger: cmd/tagger/*.go
	GOARCH=amd64 GOOS=linux go build -o bin/tagger $(PACKAGE_ROOT)/cmd/tagger

deploy:
	serverless deploy

default:
	@make compile deploy
