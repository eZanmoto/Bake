# Copyright 2012-2013 Sean Kelleher. All rights reserved.
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
TESTS=tests/perm bake bake/proj bake/recipe/test readers strio

.PHONY: all
all: clean build

# fmt before vet because fmt and vet catch the same errors, but vet outputs them
# in a way that doesn't work with Vim's Quickfix window (Vim doesn't open the
# correct file).
.PHONY: build
build: fmt vet bin/$(TARGET)

bin/$(TARGET):
	@case $$GOPATH: in \
		*/Bake:*) ;; \
		*\Bake:*) ;; \
		*) echo "No path to Bake in GOPATH" ; exit 1 ;; \
	esac
	$(GO) install $(TARGET)

.PHONY: tests
tests: build $(patsubst %,$(TSTDIR)/%.test,$(TESTS))
	echo source $(GOROOT)/src/pkg/runtime/runtime-gdb.py > $(TSTDIR)/.gdbinit

$(TSTDIR)/%.test: % $(TSTDIR)
	$(eval TST=$(subst $(SRCDIR)/,,$<))
	mkdir -p $(TSTDIR)/$(TST)
	$(GOTEST) -i $(TST)
	$(GOTEST) -c $(TST)
	mv $(notdir $(TST)).test $@

bin/rcptest: fmt vet
	$(GO) install rcptest

$(TSTDIR): $(PKGDIR)
	mkdir -p $(TSTDIR)

$(PKGDIR):
	mkdir -p $(PKGDIR)

.PHONY: runtests
runtests: build $(TESTS) testrcps
	@for TEST in $(TESTS); do \
		go test -i $$TEST; \
		go test $$TEST; \
		if [ $$? -ne 0 ]; then exit 1; fi; \
	done

.PHONY: testrcps
testrcps: bin/rcptest
	bin/rcptest

.PHONY: fmt
fmt: fmtincl fmtsrc

.PHONY: fmtincl
fmtincl:
	$(GO) install fmtincl
	./fmtincl.sh

.PHONY: fmtsrc
fmtsrc:
	gofmt -d -s $(SRCDIR)
	gofmt -s -w $(SRCDIR)

.PHONY: vet
vet:
	$(GO) tool vet $(SRCDIR)

.PHONY: clean
clean:
	rm -rf $(BINDIR) $(PKGDIR)
