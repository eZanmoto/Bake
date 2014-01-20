# Copyright {Year} {Owner}. All rights reserved.

# Targets
#
{?bin}
# all:      Make a clean build
# build:    Build the {ProjectNameLower} executable
# vet:      Runs basic safety checks on code
{?}
# clean:    Removes the local build files

.PHONY:{?bin} all build vet{?} clean
{?bin}

# Tools
GO=go
{?}

# Directories
BINDIR=bin
PKGDIR=pkg
{?bin}
SRCDIR=src

VPATH=$(SRCDIR)

# Target
TARGET={ProjectNameLower}

all: clean build

build: vet
	@case $$GOPATH: in \
		*/{ProjectName}:*) ;; \
		*\{ProjectName}:*) ;; \
		*) echo "No path to {ProjectName} in GOPATH" ; exit 1 ;; \
	esac
	$(GO) install $(TARGET)

vet:
	$(GO) tool vet $(SRCDIR)
{?}

clean:
	rm -rf $(BINDIR) $(PKGDIR)
