
# Source diagrams are all in different .pu files
SRC_DIAG=example ms-collector-status ms-management-service-components \
		 ms-service-status-changes ms-service-status ms-upload-command \
		 ms-voter-list-status

PARENT=development

# Sub-diagrams are generated from $(PARENT).pu
SUB_DIAG=processing collector_microservices collector_extension monitoring \
		 app_modules key audit


include ../../../common-model.mk


# Set environment variables for subdiagrams
collector_extension.env:
	$(eval DIAGRAM_DEF=-DCOLLECTOR_EXTENSION)

collector_microservices.env:
	$(eval DIAGRAM_DEF=-DCOLLECTOR_MICROSERVICES)

processing.env:
	$(eval DIAGRAM_DEF=-DPROCESSING_APPLICATION)

key.env:
	$(eval DIAGRAM_DEF=-DKEY_APPLICATION)

audit.env:
	$(eval DIAGRAM_DEF=-DAUDIT_APPLICATION)

monitoring.env:
	$(eval DIAGRAM_DEF=-DMONITORING)

app_modules.env:
	$(eval DIAGRAM_DEF=-DAPP_MODULES)
