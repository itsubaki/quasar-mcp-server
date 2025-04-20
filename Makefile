SHELL := /bin/bash

update:
	go get -u ./...
	go mod tidy

install:
	go install
