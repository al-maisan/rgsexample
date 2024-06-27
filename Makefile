.PHONY: all build 

BIN_DIR := ./bin
version := $(shell git rev-parse --short=12 HEAD)
timestamp := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

all: build

clean:
	rm -f $(BIN_DIR)/server
	rm -f $(BIN_DIR)/client

build: lint
	rm -f $(BIN_DIR)/server
	go build -o $(BIN_DIR)/server -v -ldflags \
		"-X main.rev=$(version) -X main.bts=$(timestamp)" cmd/server/main.go
	rm -f $(BIN_DIR)/client
	go build -o $(BIN_DIR)/client -v -ldflags \
		"-X main.rev=$(version) -X main.bts=$(timestamp)" cmd/client/main.go

lint:
	golangci-lint run
