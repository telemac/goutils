.PHONY: lint test list all

# lint the source tree
lint:
	golangci-lint run --enable-all --disable lll

# test runs tests
test:
	go test -v ./...

# list available updates in packages
list:
	go list -u -m all

all: test list lint
