SERVICES := choices mid smartid webeid proxy storage voting verification votesorder sessionstatus
OUTPUT   := $(SERVICES:%=ivxv-%@.service)

.PHONY: all
all: $(OUTPUT)

ivxv-%@.service: ivxv-$$service@.service.in
	sed 's/@service@/$*/' < '$^' > $@

	# ivxv-storage service should have a greater start timeout in case large db (> 8GB)
	# nodes will be brought up. Systemd uses TimeoutStartSec=90s by default, which may
	# be not enough.
	if [ "$*" = "storage" ]; then \
		sed -i "s/TimeoutStopSec=90s/TimeoutStopSec=90s\nTimeoutStartSec=600s/g" ivxv-$*@.service; \
	fi

.PHONY: install
install: all
	mkdir -p $(DESTDIR)/usr/lib/systemd/user
	cp $(OUTPUT) $(DESTDIR)/usr/lib/systemd/user

.PHONY: clean
clean:
	rm -f $(OUTPUT)
