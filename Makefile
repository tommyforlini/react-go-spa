GOCMD=go
GOBUILD=GOOS=linux $(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

BINARY_NAME=spa

all: clean build buildUI package

clean: 
	clear
	@echo ""
	@echo "Deleting Previous package..."
	rm -fr bin/
	rm -fr static/

build: 
	@echo ""
	@echo "Building Serverside code..."
	$(GOBUILD) -o bin/$(BINARY_NAME) -v

buildUI: 
	@echo ""
	@echo "Building Clientside code..."
	cd ui && npm run build && cd -

package:
	@echo ""
	@echo "Building Package..."
	mv static bin/
	cp .env bin/