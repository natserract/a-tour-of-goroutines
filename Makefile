-include .env

DATABASE_URL := "postgres://${DB_USERNAME}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?${DB_PARAMS}"

reset: migratedown migrateup prod
reset-dev: migratedown migrateup dev

## dev: run build and up on dev environment.
dev: build up

prod: build up-prod

## build: run build on dev environment.
build:
	GOARCH=amd64 GOOS=darwin go build -o main .

## up: run docker-compose up with dev environment.
up:
	./main

up-prod:
	ENV=production ./main

## run golang-migrate up
migrateup:
	migrate -database $(DATABASE_URL) -path db/migrations up

## run golang-migrate down
migratedown:
	migrate -database $(DATABASE_URL) -path db/migrations down
