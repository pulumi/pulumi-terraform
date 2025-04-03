MODULE          := github.com/pulumi/pulumi-terraform
VERSION         := $(shell pulumictl get version)

_ := $(shell mkdir -p bin)
_ := $(shell mkdir -p .make/sdk)
_ := $(shell go build -o bin/helpmakego github.com/iwahbe/helpmakego)

bin/pulumi-resource-terraform: $(shell bin/helpmakego provider/pulumi-resource-command)
	go build -o $@ -ldflags "-X ${MODULE}/provider/pkg/version=${VERSION}" "${MODULE}/provider/pulumi-resource-command"

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
