COMMONDIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
ROOTDIR   := $(abspath $(COMMONDIR)../../)/

include $(COMMONDIR)javavar.mk

.PHONY: all
all: GFLAGS += --warning-mode all
all:
	$(G) build installDist $(if $(DEVELOPMENT),-P development=1) $(GFLAGS)

.PHONY: all-dev
all-dev:
	$(MAKE) DEVELOPMENT=1

.PHONY: test
test:
test:
	$(G) test $(GFLAGS)

.PHONY: install
install: GFLAGS += -x test
install: all
ifndef DESTDIR
	$(error $$DESTDIR must be set to install Java applications)
endif
	if [ -d build/distributions ]; then \
		mkdir -p $(DESTDIR)/ && \
		cp build/distributions/*.zip $(DESTDIR)/; \
	fi

.PHONY: clean
clean:
clean:
	$(G) clean $(GFLAGS)
