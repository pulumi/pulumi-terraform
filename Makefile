MODULE          := github.com/pulumi/pulumi-terraform
VERSION         := $(shell pulumictl get version)

.PHONY: all
all: schema.json build_sdks bin/pulumi-resource-terraform

_ := $(shell mkdir -p bin)
_ := $(shell mkdir -p .make/sdk)
_ := $(shell go build -o bin/helpmakego github.com/iwahbe/helpmakego)

bin/pulumi-resource-terraform: $(shell bin/helpmakego .)
	go build -o $@ -ldflags "-X ${MODULE}/provider/version.Version=${VERSION}" "${MODULE}"

schema.json: bin/pulumi-resource-terraform
	pulumi package get-schema $< > $@

.PHONY: .make/phony/sdk/%
.make/phony/sdk/%: bin/pulumi-resource-terraform
	pulumi package gen-sdk $< --language $*

.PHONY: build_go build_nodejs build_python build_java build_dotnet build_sdks

build_sdks: build_go build_nodejs build_python build_java build_dotnet

build_go:     .make/phony/sdk/go
build_nodejs: .make/phony/sdk/nodejs
build_python: .make/phony/sdk/python
build_java:   .make/phony/sdk/java
build_dotnet: .make/phony/sdk/dotnet

lint:
	golangci-lint run --config ./.golangci.yml

.PHONY: test test_unit test_integration

test: test_unit test_integration

test_unit:
	go test $$(go list ./... | grep -v /examples)

# By default, `$(MAKE) test_integration` will run all integration tests.
#
# To run a specific integration test, you can override TAGS:
#
#     make test_integration TAGS=yaml
#
test_integration: TAGS ?= all
test_integration: bin/pulumi-resource-terraform
	go test $$(go list ./... | grep /examples) -tags ${TAGS} -count 1
