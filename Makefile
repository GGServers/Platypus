.DEFAULT_GOAL := build
PLATYPUS_VERSION := 2.0.0
# Workaround for CircleCI Docker image and mkdir.
SHELL := /bin/bash

make_build_dir:
	mkdir -p build/{bin,assets,tars}

build: make_build_dir
	go build -o build/bin/platypus

small: make_build_dir
	go build -o build/bin/platypus -ldflags="-s -w"
	upx --brute build/bin/platypus -9 --no-progress

client: make_build_dir
	go build -o build/bin/platypus-client -ldflags="-s -w" client/client.go
	upx --brute build/bin/platypus-client -9 --no-progress

client_pi: make_build_dir
	env GOOS=linux GOARCH=arm GOARM=5 go build -o build/bin/platypus-client-arm -ldflags="-s -w" client/client.go 
	upx --brute build/bin/platypus-client-arm -9 --no-progress

run:
	go run main.go

dist: clean make_build_dir small client client_pi

clean:
	rm -rf build
