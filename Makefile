#
# By default this will build the project on every non-mobile platform
# supported by the installed go environment.
#
# To limit a build to a single environment, you can force it to just a
# single platform by prefixing make with:
#
# PLATFORMS=linux:amd64: make clean all
#
# Just change the entry for your OS and CPU. These are listed in platforms.md
#
# Note: For 32 bit arm processors the 3rd parameter is important.
# e.g. use linux:arm:6 or linux:arm:7
#
# For all other processors, including arm64, leave the third field blank.
#
# To disable tests, you can prefix make with:
#
# GO_TEST="#" make clean all
#
# The quotes are important!
#
# You can combine the two as necessary.
#
# e.g. GO_TEST="#" PLATFORMS=linux:amd64: make clean all
#
# For a parallel builds you can use the -j parameter to make as usual.
#
# e.g.: make -j 8 clean all
#
# Pick a value suitable to the number of cores/thread your machine has.
# This is useful for a full build of all platforms as it will build all
# of the binaries in parallel speeding up the full build.
#

# The repository name/package prefix.
# This should match the value of module in go.mod
PACKAGE_PREFIX = $(shell grep ^module go.mod | cut -f2 -d' ' | head -1)
PACKAGE_NAME = $(shell basename $(PACKAGE_PREFIX))
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null | sed "s/-/./g")
DIST_PREFIX = $(PACKAGE_NAME)_$(VERSION)
BUILD_DATE = $(shell date)

DOCKER-REPOSITORY = docker.europa.area51.dev/
DOCKER-TAG = $(DOCKER-REPOSITORY)$(PACKAGE_NAME):$(if $(filter dirty,$(VERSION)),"latest",$(VERSION))

# Where to place build artifacts. These must be subdirectories here and not
# a path elsewhere, otherwise it will break the build!
BUILDS 	= builds
DIST		= dist

# BINDIR is the prefix before any built tools. Set to "" for none, otherwise
# it must end with /
BINDIR ?= bin/

.PHONY: all clean init test tools dist real-clean

all: init test tools

include Makefile.include
include Go.include
include Docker.include

clean: $(targets-clean)
	$(call REMOVE,$(BUILDS) $(DIST))

real-clean: clean $(targets-real-clean)

init: $(targets-validate) $(resolve-platforms) $(targets-init)
	@$(ECHO) "$(PLATFORMS)"

test: $(targets-test)

tools: $(targets-tools)

dist1: all
	$(MKDIR) $(DIST)
	$(foreach PLATFORM,$(shell cd $(BUILDS);ls -d */*),$(call TAR,$(PLATFORM))${\n})

dist: all
	$(foreach PLATFORM,$(shell cd $(BUILDS);ls -d */*),@$(MAKE) $(DIST)/$(call TAR-FILE,$(PLATFORM)){\n})
	@$(MAKE) $(targets-dist)

$(DIST)/$(PACKAGE_NAME)_$(VERSION)-%.tgz: builds/$(subst _,/,%)
	$(MKDIR) $(DIST)
	echo "$@ $<"
