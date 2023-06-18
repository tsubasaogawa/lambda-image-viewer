.PHONY: compile deploy default
.DEFAULT_GOAL := default

compile: *.go
	GOARCH=amd64 GOOS=linux go build -o bin/main

deploy:
	serverless deploy

default:
	@make compile deploy
