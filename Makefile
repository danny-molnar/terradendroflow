# Makefile for TDFlow

.PHONY: test lint build clean

test:
	go test ./... -v

lint:
	golangci-lint run

build:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/tdflow_linux main.go
	GOOS=darwin GOARCH=amd64 go build -o release/tdflow_mac main.go
	GOOS=windows GOARCH=amd64 go build -o release/tdflow_windows.exe main.go

clean:
	rm -rf release
