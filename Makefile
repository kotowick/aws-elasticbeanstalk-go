
# Generate tarball with new build of osx_vpn_manager
#
# NOTE: OSX only
VERSION=$$(cat main.go | grep -i "cliVersion =" | awk {'print$$3'} | tr -d '"')


all: clean build compress report

clean:
	@rm -f /tmp/ebd-*.tar.gz
	@rm -f ./bin/ebd

build:
	@echo Building ebd version $(VERSION)
	@go build -o ./bin/ebd

compress:
	@tar czf /tmp/ebd-$(VERSION).tar.gz ./versions

report:
	@rm -f ./bin/ebd
	@shasum -a 256 /tmp/ebd-$(VERSION).tar.gz

.PHONY: all clean build
