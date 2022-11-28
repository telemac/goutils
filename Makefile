.PHONY: lint test list vulncheck all

# lint the source tree
lint:
	CGOENAGLE=0 GODEBUG=netdns=go golangci-lint run --enable-all --disable lll

# test runs tests
test:
	go clean -testcache
	CGOENAGLE=0 GODEBUG=netdns=go go test -race ./...

# list available updates in packages
list:
	go list -u -m all

# vulnerability check
vulncheck:
	govulncheck ./...

all: test list lint
