TEST?=./...
NAME=aws-ssh
# NAME = $(shell awk -F\" '/^const Name/ { print $$2 }' main.go)
VERSION = $(shell awk -F\" '/^const Version/ { print $$2 }' cmd/versionCmd.go)

all: build

build: deps
	@mkdir -p bin/
	@echo ${VERSION}
	go build -o bin/$(NAME)

xcompile: deps
	@rm -rf build/
	@mkdir -p build
	gox \
    -arch="amd64" \
		-os="darwin" \
		-os="linux" \
		-os="windows" \
		-output="build/{{.Dir}}_$(VERSION)_{{.OS}}_{{.Arch}}/$(NAME)"

package: xcompile
	$(eval FILES := $(shell ls build))
	@mkdir -p build/tgz
	for f in $(FILES); do \
		(cd $(shell pwd)/build && tar -zcvf tgz/$$f.tar.gz $$f); \
		echo $$f; \
	done

.PHONY: all deps build xcompile package
