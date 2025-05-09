# Build variables
BINARY_NAME = dasung_auto
BUILD_DIR = build

# Platform specific variables
DARWIN = darwin

# Architecture specific variables
# AMD64 = amd64
ARM64 = arm64

# Build flags
LDFLAGS = -ldflags "-s -w"

# Default target
.PHONY: all
all: clean build

# Build all platforms
.PHONY: build
build: darwin

# Darwin (macOS) builds
.PHONY: darwin
darwin: darwin-arm64 #darwin-amd64 

# .PHONY: darwin-amd64
# darwin-amd64:
# 	GOOS=$(DARWIN) GOARCH=$(AMD64) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(DARWIN)-$(AMD64)

.PHONY: darwin-arm64
darwin-arm64:
	GOOS=$(DARWIN) GOARCH=$(ARM64) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(DARWIN)-$(ARM64)

# Clean build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Build all platforms (default)"
	@echo "  build        - Build all platforms"
	@echo "  darwin       - Build all macOS versions"
	@echo "  clean        - Clean build directory"
	@echo "  help         - Show this help message"
