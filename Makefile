#
#  As a quick totorial for Makefile check https://makefiletutorial.com/
#

.DEFAULT_GOAL := build

binary_name = user-service

#
# Install all required for build dependencies
#
preinstall:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.48.0
.PHONY:preinstall

#
# Format all files in-place.
# Check: https://go.dev/blog/gofmt
#
fmt:
	gofmt -w  .
.PHONY:fmt

#
# Run static analysis (aka linting)
# Check: https://golangci-lint.run/
#
lint:	
	golangci-lint run 
.PHONY:lint

build: fmt lint
	go build -o bin/$(binary_name)
.PHONY:build

run: build
	bin/$(binary_name)
.PHONY:run

#
# cross-compilation binaries for Win, MacOS and Linux
#
release: fmt lint
	GOOS=windows GOARCH=amd64 go build -o bin/$(binary_name)_windows_amd64
	GOOS=linux GOARCH=amd64 go build -o bin/$(binary_name)_linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o bin/$(binary_name)_darwin_amd64
.PHONY:release

#
# Clean up all artifacts
#
clean:
	rm -rf bin/*
.PHONY:clean

