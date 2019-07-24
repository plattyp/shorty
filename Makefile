# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

build:
	$(GOCLEAN)
	$(GOBUILD) -i -o bin/shorty main.go

install:
	createdb shorty
	./bin/migrate up

run:
	make build
	bin/shorty

migrate:
	./bin/migrate up

migrate-deployed:
	./bin/migrate up --env deployed

init-test:
	dropdb shorty-test --if-exists
	createdb shorty-test
	./bin/migrate up --env test

init-test-travis:
	dropdb shorty-test --if-exists
	createdb shorty-test
	./bin/migrate up --env test

test:
	go test
