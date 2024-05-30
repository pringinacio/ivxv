# Common Makefile for all Go modules
# 1. Create new Makefile in a Go module directory
# 2. Write as first line in a Makefile `include ../common/go/common.mk`

COMMONDIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
ROOTDIR   := $(abspath $(COMMONDIR)../../)/

include $(COMMONDIR)govar.mk

.PHONY: all
all: tidy generate
	env GOBIN=$(or $(IVXV_GOBIN),$(abspath bin)) \
		$(GO) install $(if $(DEVELOPMENT),-tags development )./...

# Generate gen_types.go for each package
.PHONY: gotools
gotools:
	$(MAKE) -C $(ROOTDIR)common/tools/go all ROOTDIR=$(ROOTDIR)

.PHONY: lint
lint: all
	if which golangci-lint > /dev/null; then \
		env PATH="$(dir $(GO)):$$PATH" \
			golangci-lint run --config $(COMMONDIR)golangci-lint.yaml; \
	fi

.PHONY: test
test: lint
	env GODEBUG=x509sha1=1 $(GO) test -v $(GOTESTFLAGS) ./...

.PHONY: install
install: all
	if [ -d bin ]; then mkdir -p $(DESTDIR)/usr && cp -r bin $(DESTDIR)/usr; fi
	$(if $(EXTRADATA),mkdir -p $(DESTDIR)/usr/share/ivxv && cp -r $(EXTRADATA) $(DESTDIR)/usr/share/ivxv)

.PHONY: generate
generate: gotools
	$(GO) generate ./...
	$(ROOTDIR)common/tools/go/bin/gen -base $(ROOTDIR) ./...

.PHONY: clean
clean:
	find . -not -path \*/testdata/\* \( \
		-name gen_types.go -o \
		-name gen_types_test.go -o \
		-name gen_import.go -o \
		-name gen_import_dev.go \
		\) -delete
	 -$(GO) clean -i ./...
	rm -rf bin/ pkg/

.PHONY: tidy
tidy:
	$(GO) mod tidy

# Generate gen_import.go for each package where in source files:
# `import //ivxv:modules <module>`
# Generate gen_import_dev.go for each package where in source files:
# `import //ivxv:development <module>`
.PHONY: goimports
goimports:
	$(ROOTDIR)scripts/goimports.sh ./...
