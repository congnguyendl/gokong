.DEFAULT_GOAL := default

TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

ver ?= 6.0.3

update:
	@if git rev-parse v$(ver) >/dev/null 2>&1; then \
		echo "Tag v$(ver) already exists. Update the version."; \
		exit 1; \
	fi
	@git add .
	@git commit -m "update: sdk v$(ver)" || echo "No changes to commit"
	@git tag -a v$(ver) -m "update: sdk v$(ver)"
	@git push origin master
	@git push origin --tags


default: build test

build: fmtcheck errcheck vet
	go install

test: goimportscheck
	go test -v --tags="community" ./...

test-all: goimportscheck
	go test -v --tags="all" ./...

test-enterprise: goimportscheck
	go test -v --tags="enterprise" ./...

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

goimports:
	goimports -w $(GOFMT_FILES)

goimportscheck:
	@sh -c "'$(CURDIR)/scripts/goimportscheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

install-goimports:
	@go get golang.org/x/tools/cmd/goimports

vendor-status:
	@govendor status

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./aws"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: build test testacc vet fmt fmtcheck errcheck vendor-status test-compile
