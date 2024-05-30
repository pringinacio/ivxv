# common.mk: Common recipes for Documentation.
#
# Defines a catch-all recipe which routes unknown targets to Sphinx in "make
# mode". If $(DEPENDENCIES) is defined, then it will be set as a prerequisite
# of the catch-all and removed on clean. In addition, if the source directory
# contains a "model" subdirectory, then it invokes make in it to build and
# clean that too.
#
# Defines a "diff" recipe which generates a PDF that highlights differences
# between the current and latest released version of the document. If there are
# no differences, then no PDF is produced.
#
# Defines "install-pdf", "install-html", and "install-diff" recipes which build
# the specified type of documentation and install it to $(DESTDIR).

common.mk := $(lastword $(MAKEFILE_LIST))
Makefile  := $(lastword $(filter-out $(common.mk),$(MAKEFILE_LIST)))

# You can set these variables from the command line, and also
# from the environment for the first two.
SPHINXOPTS    ?= -c $(dir $(common.mk))
SPHINXBUILD   ?= sphinx-build
SOURCEDIR     = $(dir $(Makefile))
BUILDDIR      = $(SOURCEDIR)_build
SPHINXINTL    ?= sphinx-intl

# Set IVXV_DOCUMENT to the name of the source directory. This will be used as
# the key to look up configuration from documents.py.
export IVXV_DOCUMENT := $(notdir $(patsubst %/,%,$(abspath $(SOURCEDIR))))

# Special case the help target: set as default and skip prerequisites.
.DEFAULT_GOAL := help
.PHONY: help
help:
	@$(SPHINXBUILD) -M help "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)


translation:
	@$(SPHINXBUILD) -M gettext "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)
	@$(SPHINXINTL) update -p _build/gettext -l en


english: clean
	export SPHINXOPTS="-D language='en'" && $(SPHINXBUILD) -M latexpdf "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

estonian: clean
	@$(SPHINXBUILD) -M latexpdf "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)


# Special case the clean target: skip prerequisites and perform extra steps.
.PHONY: clean
clean:
	#if [ -d "model" ]; then $(MAKE) -C model clean; fi
	@$(SPHINXBUILD) -M clean "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)
	rm -rf $(BUILDDIR) $(DEPENDENCIES)

# Do not regenerate the Makefiles.
$(common.mk) $(Makefile): ;

# Catch-all target: route all unknown targets to Sphinx using the new
# "make mode" option.  $(O) is meant as a shortcut for $(SPHINXOPTS).
%: $(DEPENDENCIES)
	if [ -d "model" ]; then $(MAKE) -C model; fi
	$(SPHINXBUILD) -M $@ "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

# Variables and rules for diffing.
differ.sh := $(dir $(common.mk))differ.sh
master    := $(dir $(common.mk))_master/
gitdir    := $(patsubst $(shell git rev-parse --show-toplevel)%,%,$(abspath $(SOURCEDIR)))
olddir    := $(master)$(gitdir)

ifeq ("$(olddir)", "../../_master//Documentation/public/protokollid")
	olddir := "../../_master//Documentation/et/protokollid"
endif

ifeq ("$(olddir)", "../../_master//Documentation/public/uldsisukord")
	olddir := "../../_master//Documentation/et/uldsisukord"
endif

ifeq ("$(olddir)", "../../_master//Documentation/public/arhitektuur")
	olddir := "../../_master//Documentation/et/arhitektuur"
endif

.PHONY: diff
diff: latex master-latex
	# latexmk attempts to proceed even though errors may be present
	# pdflatex attempts to proceed with errors and is silent (batchmode)
	LATEXMKOPTS="-f -interaction=batchmode" $(differ.sh) $(BUILDDIR)/master $(BUILDDIR)

.PHONY: master-latex
master-latex: $(master)
	if [ -d $(olddir) ]; then \
		$(MAKE) -C $(olddir) BUILDDIR=$(abspath $(BUILDDIR))/master latex; \
	fi

$(master):
	git worktree add $@ 1.8.3

# Installation rules.
.PHONY: install-pdf
install-pdf: estonian
	cp --update $(filter-out %-diff.pdf,$(wildcard $(BUILDDIR)/latex/*.pdf)) "$(DESTDIR)"

.PHONY: install-en-pdf
install-en-pdf: english
	cp --update $(filter-out %-diff.pdf,$(wildcard $(BUILDDIR)/latex/*.pdf)) "$(DESTDIR)"

.PHONY: install-html
install-html: html
	cp --recursive --update $(BUILDDIR)/html/** "$(DESTDIR)"

.PHONY: install-diff
install-diff: diff
	$(eval diff.pdf = $(wildcard $(BUILDDIR)/latex/*-diff.pdf))
	$(if $(diff.pdf),cp --update $(diff.pdf) "$(DESTDIR)")
