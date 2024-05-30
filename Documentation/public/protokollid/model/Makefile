# Source diagrams are all in different .pu files
SRC_DIAG=

PARENT=sequence

# Sub-diagrams are generated from $(PARENT).pu
SUB_DIAG=phase1 phase2 phase3


include ../../../common-model.mk


# Set environment variables for subdiagrams
phase1.env:
	$(eval DIAGRAM_DEF=-DPHASE_1)

phase2.env:
	$(eval DIAGRAM_DEF=-DPHASE_2)

phase3.env:
	$(eval DIAGRAM_DEF=-DPHASE_3)
