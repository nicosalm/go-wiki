# Go parameters
GO = go
GOFILES = wiki.go

# Targets
.PHONY: all run clean

# Default target
all: run

# Build and run the Go application
run:
	$(GO) run $(GOFILES)

# Clean up binary files
clean:
	$(GO) clean
	rm -f *.out

