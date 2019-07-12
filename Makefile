# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

build:
	$(GOCLEAN)
	$(GOBUILD) -i -o bin/shorty main.go

run:
	make build
	bin/shorty

migrate:
	sql-migrate up