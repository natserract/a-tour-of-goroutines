-include .env

## dev: run build and up on dev environment.
dev: build up

## build: run build on dev environment.
build:
	GOARCH=amd64 GOOS=darwin go build -o main .

## up: run docker-compose up with dev environment.
up:
	./main
