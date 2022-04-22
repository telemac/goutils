.PHONY: lint test list all

# lint the source tree
lint:
	CGOENAGLE=0 GODEBUG=netdns=go golangci-lint run --enable-all --disable lll

# test runs tests
test:
	CGOENAGLE=0 GODEBUG=netdns=go go test -race ./...

# list available updates in packages
list:
	go list -u -m all

all: test list lint
