# Copyright 2012 Sean Kelleher. All rights reserved.
# Use of this source code is governed by a GPL
# license that can be found in the LICENSE file.

# Targets
#
# all:      Make a clean build
# build:    Build the bake executable
# tests:    Build the debuggable testing files
# runtests: Runs all tests
# fmt:      Checks code and resource formatting
# vet:      Runs basic safety checks on code
# clean:    Removes the local build files

.PHONY: all build tests fmt fmtincl fmtsrc runtests vet clean

# Tools
GO=go
GOTEST=$(GO) test

# Directories
BINDIR=bin
PKGDIR=pkg
SRCDIR=src

VPATH=$(SRCDIR)

# Target
TARGET=bake

# Testing
TSTDIR=$(PKGDIR)/tst
TESTS=tests/perm bake bake/proj

all: clean build

build: vet fmt
	$(GO) install $(TARGET)

tests: build $(patsubst %,$(TSTDIR)/%.test,$(TESTS))
	echo source $(GOROOT)/src/pkg/runtime/runtime-gdb.py > $(TSTDIR)/.gdbinit

$(TSTDIR)/%.test: % $(TSTDIR)
	$(eval TST=$(subst $(SRCDIR)/,,$<))
	mkdir -p $(TSTDIR)/$(TST)
	$(GOTEST) -i $(TST)
	$(GOTEST) -c $(TST)
	mv $(notdir $(TST)).test $@

$(TSTDIR): $(PKGDIR)
	mkdir -p $(TSTDIR)

$(PKGDIR):
	mkdir -p $(PKGDIR)

runtests: build $(TESTS)
	@for TEST in $(TESTS); do \
		go test $$TEST; \
	done

fmt: fmtincl

fmtincl: fmtsrc
	$(GO) install fmtincl
	./fmtincl.sh

fmtsrc:
	./fmtsrc.sh

vet:
	$(GO) tool vet $(SRCDIR)

clean:
	rm -rf $(BINDIR) $(PKGDIR)
