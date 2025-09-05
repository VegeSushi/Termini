BINARY_NAME = termini
SRC = main.go
INSTALL_DIR = /usr/local/bin

GOOS ?= linux
GOARCH ?= amd64

.PHONY: all
all: build

.PHONY: build
build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINARY_NAME) $(SRC)
	@echo "Build complete!"

.PHONY: install
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	sudo cp $(BINARY_NAME) $(INSTALL_DIR)/
	sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installation complete!"

.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -f $(BINARY_NAME)
	@echo "Done."

.PHONY: uninstall
uninstall:
	@echo "Removing $(INSTALL_DIR)/$(BINARY_NAME)..."
	sudo rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Uninstalled."
