PROJECT_NAME := Pulumi Terraform Resource Provider
include build/common.mk

PACK             := terraform
PACKDIR          := sdk
NODE_MODULE_NAME := @pulumi/terraform
NUGET_PKG_NAME   := Pulumi.Terraform
PROJECT          := github.com/pulumi/pulumi-terraform
TESTPARALLELISM  := 10

VERSION          ?= $(shell scripts/get-version)
PYPI_VERSION     := $(shell scripts/get-py-version)

VERSION_FLAGS    := -ldflags "-X github.com/pulumi/pulumi-terraform/provider/v2/pkg/version.Version=${VERSION}"

DOTNET_PREFIX  := $(firstword $(subst -, ,${VERSION:v%=%})) # e.g. 1.5.0
DOTNET_SUFFIX  := $(word 2,$(subst -, ,${VERSION:v%=%}))    # e.g. alpha.1

ifeq ($(strip ${DOTNET_SUFFIX}),)
	DOTNET_VERSION := $(strip ${DOTNET_PREFIX})
else
	DOTNET_VERSION := $(strip ${DOTNET_PREFIX})-$(strip ${DOTNET_SUFFIX})
endif

build::
	cd provider && go install $(VERSION_FLAGS) ${PROJECT}/provider/v2/cmd/pulumi-resource-${PACK}
	cd ${PACKDIR}/nodejs/ && \
		yarn install && \
		yarn run tsc
	cp LICENSE ${PACKDIR}/nodejs/package.json ${PACKDIR}/nodejs/yarn.lock \
		${PACKDIR}/nodejs/bin
	cp README.md ${PACKDIR}/nodejs/bin/README.md
	cp README.md ${PACKDIR}/python/README.md
	sed -i.bak 's/$${VERSION}/$(VERSION)/g' ${PACKDIR}/nodejs/bin/package.json
	cd ${PACKDIR}/python/ && \
		$(PYTHON) setup.py clean --all 2>/dev/null && \
		rm -rf ./bin/ ../python.bin/ && cp -R . ../python.bin && mv ../python.bin ./bin && \
		sed -i.bak -e "s/\$${VERSION}/$(PYPI_VERSION)/g" -e "s/\$${PLUGIN_VERSION}/$(VERSION)/g" ./bin/setup.py && \
		cd ./bin && $(PYTHON) setup.py build sdist
	cd ${PACKDIR}/dotnet/ && \
		echo "${VERSION:v%=%}" >version.txt && \
		dotnet build /p:Version=${DOTNET_VERSION}

lint::
	#cd provider/cmd/pulumi-resource-terraform && golangci-lint run

install:: build
	[ ! -e "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)" ] || rm -rf "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)"
	mkdir -p "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)"
	cp -r sdk/nodejs/bin/. "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)"
	rm -rf "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)/node_modules"
	rm -rf "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)/tests"
	cd "$(PULUMI_NODE_MODULES)/$(NODE_MODULE_NAME)" && \
		yarn install --offline --production && \
		(yarn unlink > /dev/null 2>&1 || true) && \
		yarn link
	[ ! -e "$(PULUMI_NUGET)" ] || rm -rf "$(PULUMI_NUGET)/*"
	find . -name '$(NUGET_PKG_NAME).*.nupkg' -exec cp -p {} ${PULUMI_NUGET} \;

test_fast::
	cd examples && $(GO_TEST_FAST) .

test_all::
	cd examples && $(GO_TEST) .

.PHONY: publish_tgz
publish_tgz:
	$(call STEP_MESSAGE)
	./scripts/publish_tgz.sh

.PHONY: check_clean_worktree
check_clean_worktree:
	$$(go env GOPATH)/src/github.com/pulumi/scripts/ci/check-worktree-is-clean.sh

# While pulumi-terraform is not built using tfgen/tfbridge, the layout of the source tree is the same as these
# style of repositories, so we can re-use the common publishing scripts.
.PHONY: publish_packages
publish_packages:
	$(call STEP_MESSAGE)
	$$(go env GOPATH)/src/github.com/pulumi/scripts/ci/publish-tfgen-package .
	$$(go env GOPATH)/src/github.com/pulumi/scripts/ci/build-package-docs.sh terraform

# The travis_* targets are entrypoints for CI.
.PHONY: travis_cron travis_push travis_pull_request travis_api
travis_cron: all
travis_push: all check_clean_worktree only_test publish_tgz publish_packages
travis_pull_request: all
travis_api: all
