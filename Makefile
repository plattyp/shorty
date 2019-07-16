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
	sql-migrate up

run:
	make build
	bin/shorty

migrate:
	sql-migrate up

init-test:
	dropdb shorty-test --if-exists
	createdb shorty-test
	sql-migrate up --env test

init-test-travis:
	./sql-migrate up --env travis

test:
	go test