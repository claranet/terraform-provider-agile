GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
ACCTEST_PARALLELISM ?= 20

ifneq ($(origin TESTS), undefined)
	RUNARGS = -run='$(TESTS)'
endif

default: build

# Build
build: fmtcheck
	go install

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

# Run acceptance tests
testacc: fmtcheck
	TF_ACC=1 go test ./... -v -parallel $(ACCTEST_PARALLELISM) $(RUNARGS) $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

generate-docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

# Deploy Local
build-local: build
	mv ${GOPATH}/bin/terraform-provider-agile ~/.terraform.d/plugin-cache/registry.terraform.io/claranet/agile/${TERRAFORM_AGILE_VERSION}/linux_amd64/terraform-provider-agile_v${TERRAFORM_AGILE_VERSION}

.PHONY: build fmtcheck testacc build-local vet fmt errcheck generate-docs