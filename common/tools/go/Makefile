include ../../../common/go/govar.mk

.PHONY: all
all: bin/gen

bin/gen:
	$(GO) generate ./...
	env GOBIN=$(or $(IVXV_GOBIN),$(abspath bin)) \
		$(GO) install $(if $(DEVELOPMENT),-tags development )./...

.PHONY: clean
clean:
	rm -rf cmd/gen/gentmpl.go
