SHELL=/bin/bash

GOTOOLCHAIN=go1.24.4

PNAME = golang-task-api

.PHONY: clear prep perm

all: prep static perm

clear:
	rm -rf bin/*

prep:
	GOTOOLCHAIN=$(GOTOOLCHAIN) go mod tidy

static:
	GOOS=linux GOARCH=amd64 GOTOOLCHAIN=$(GOTOOLCHAIN) go build -buildvcs=false -o=bin/$(PNAME)-linux-amd64 -o=bin/$(PNAME) ./cmd/...
	GOOS=windows GOARCH=amd64 GOTOOLCHAIN=$(GOTOOLCHAIN) go build -buildvcs=false -o=bin/$(PNAME)-windows-amd64.exe ./cmd/...

perm:
	chmod -R +x bin