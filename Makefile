MODULE          := github.com/pulumi/pulumi-terraform
VERSION         := 6.0.0 # $(shell pulumictl get version)

.PHONY: all
all: schema.json build_sdks bin/pulumi-resource-terraform

_ := $(shell mkdir -p bin)
_ := $(shell mkdir -p .make/sdk)
_ := $(shell go build -o bin/helpmakego github.com/iwahbe/helpmakego)

bin/pulumi-resource-terraform: $(shell bin/helpmakego .)
	go build -o $@ -ldflags "-X ${MODULE}/provider/version.version=${VERSION}" "${MODULE}"

schema.json: bin/pulumi-resource-terraform
	pulumi package get-schema $< > $@

.PHONY: .make/phony/sdk/%
.make/phony/sdk/%: bin/pulumi-resource-terraform
	pulumi package gen-sdk $< --language $*

.PHONY: build_go build_nodejs build_python build_java build_dotnet build_sdks

build_sdks: build_go build_nodejs build_python build_java build_dotnet

build_go:     .make/phony/sdk/go
build_nodejs: .make/phony/sdk/nodejs
	cd sdk/nodejs && yarn install && yarn run tsc
	cp README.md LICENSE sdk/nodejs/package.json sdk/nodejs/yarn.lock sdk/nodejs/bin/

build_python: .make/phony/sdk/python
build_java:   .make/phony/sdk/java
build_dotnet: .make/phony/sdk/dotnet

lint:
	golangci-lint run --config ./.golangci.yml --build-tags all

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
	go test $$(go list ./... | grep /examples) -tags ${TAGS} -count 1 -v

# To make an immediately observable change to .ci-mgmt.yaml:
#
# - Edit .ci-mgmt.yaml
# - Run make ci-mgmt to apply the change locally.
#
ci-mgmt: .ci-mgmt.yaml
	go run github.com/pulumi/ci-mgmt/provider-ci@b6bfde4bf3d1f9e539671e20aad7801e4ba5d300 generate
.PHONY: ci-mgmt

# Targets for ci-mgmt (also includes the build_% category of commands)
.PHONY: codegen generate_schema local_generate provider test_provider \
	install_go_sdk install_nodejs_sdk install_python_sdk install_java_sdk install_dotnet_sdk \
	generate_go generate_nodejs generate_python generate_java generate_dotnet

codegen: schema.json build_sdks
generate_schema: schema.json
local_generate: # It's not clear what this should do
install_go_sdk:
	# "This is a no-op that satisfies ci-mgmt
install_nodejs_sdk: build_nodejs
	-yarn unlink --cwd sdk/nodejs/bin
	yarn link --cwd sdk/nodejs/bin
install_python_sdk:
	# "This is a no-op that satisfies ci-mgmt
install_java_sdk:
	# "This is a no-op that satisfies ci-mgmt
install_dotnet_sdk:
	# "This is a no-op that satisfies ci-mgmt
generate_go: build_go
generate_nodejs: build_nodejs
generate_python: build_python
generate_java: build_java
generate_dotnet: build_dotnet
provider: bin/pulumi-resource-terraform
test_provider: test_unit
