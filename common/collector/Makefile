.DEFAULT_GOAL := all

.PHONY: install
install: all
	$(MAKE) -C config $@
	$(MAKE) -C scripts $@
	$(MAKE) -f ../go/common.mk $@

# Remove additional files in clean.
CLEAN := common/collector/server/tlsciphersuites.go

.PHONY: clean
clean:
	$(MAKE) -f ../go/common.mk clean
	rm -rf $(CLEAN)

%:
	$(MAKE) -f ../go/common.mk $@
