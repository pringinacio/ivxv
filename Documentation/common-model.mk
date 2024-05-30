PLANTUML=env -u DISPLAY plantuml
PLANTUML_EN=env -u DISPLAY plantuml -I./en/lang.pu
PLANTUML_ET=env -u DISPLAY plantuml -I./et/lang.pu

LANGUAGES=et en

DIAGRAMS=$(SRC_DIAG) $(SUB_DIAG)

# We process each diagram in all languages, so that the variables are set properly
TARGETS=$(foreach DIAGRAM,$(DIAGRAMS),$(patsubst %, img/$(DIAGRAM).%.png, $(LANGUAGES)))

all:
	@echo "Dry run, make model to apply changes"
	$(MAKE) -n model

model: $(TARGETS)

# Process manually translated diagrams
img/%.png: %.pu
	$(PLANTUML) $<
	mv $*.png $@


# Process source diagrams
img/%.et.png: %.pu
	$(PLANTUML_ET) $<
	mv $*.png $@

img/%.en.png: %.pu
	$(PLANTUML_EN) $<
	mv $*.png $@


# Process sub-diagrams of major diagram
img/%.et.png: $(PARENT).pu %.env
	$(PLANTUML_ET) $(DIAGRAM_DEF) $<
	mv $(PARENT).png $@

img/%.en.png: $(PARENT).pu %.env
	$(PLANTUML_EN) $(DIAGRAM_DEF) $<
	mv $(PARENT).png $@
