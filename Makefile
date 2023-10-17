.PHONY: *

BINARY_NAME := vault-list

build:
	go build -o bin/$(BINARY_NAME)

install: build
	cp bin/$(BINARY_NAME) /usr/local/bin/
