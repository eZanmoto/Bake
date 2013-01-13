# Copyright {Year} {Owner}. All rights reserved.

# Targets
#
{?bin}
# all:      Make a clean build
# build:    Build the {ProjectNameLower} executable
{!}
{?test}
# tests:    Build the debuggable testing files
# runtests: Runs all tests
{!}
{?bin}
# vet:      Runs basic safety checks on code
{!}
# clean:    Removes the local build files

.PHONY:{?bin} all build{!}{?test} tests runtests{!}{?bin} vet{!} clean

# Tools
GO=go
{?test}
GOTEST=$(GO) test
{!}

# Directories
BINDIR=bin
PKGDIR=pkg
{?bin}
SRCDIR=src

VPATH=$(SRCDIR)

# Target
TARGET={ProjectNameLower}
{!}
{?test}

# Testing
TSTDIR=$(PKGDIR)/tst
TESTS={ProjectNameLower}
{!}
{?bin}

all: clean build

build: vet
	@case $$GOPATH: in \
		*/{ProjectName}:*) ;; \
		*) echo "No path to {ProjectName} in GOPATH" ; exit 1 ;; \
	esac
	$(GO) install $(TARGET)
{!}
{?test}

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
{!}
{?bin}

vet:
	$(GO) tool vet $(SRCDIR)
{!}

clean:
	rm -rf $(BINDIR) $(PKGDIR)
