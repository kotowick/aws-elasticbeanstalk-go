
# Generate tarball with new build of osx_vpn_manager
#
# NOTE: OSX only
VERSION=$$(cat main.go | grep -i "cliVersion =" | awk {'print$$3'} | tr -d '"')
BINARY_NAME=ebd

all: clean build compress report

clean:
	@rm -f /tmp/ebd-*.tar.gz
	@rm -f ./bin/ebd

build:
	@echo Building $(BINARY_NAME) version $(VERSION)
	@go build -a -tags netgo -ldflags '-w' -o ./bin/$(BINARY_NAME)-$(VERSION)
	@cp ./bin/$(BINARY_NAME)-$(VERSION) ./bin/$(BINARY_NAME)-latest

compress:
	@tar czf /tmp/$(BINARY_NAME)-$(VERSION).tar.gz ./versions

report:
	@rm -f ./bin/$(BINARY_NAME)
	@shasum -a 256 /tmp/$(BINARY_NAME)-$(VERSION).tar.gz

.PHONY: all clean build
