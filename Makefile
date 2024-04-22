# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Directories
SRC_DIR=./
TEST_DIR=./
BENCH_DIR=./

# Targets
all: test benchmark

test:
	$(GOTEST) -v $(SRC_DIR)/...

benchmark:
	$(GOTEST) -v -bench=. -benchmem $(BENCH_DIR)/...

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

.PHONY: all test benchmark clean