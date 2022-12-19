BINARY=remote-access

GOOS=linux GOARCH=amd64 go build -trimpath -o ./bin/linux/amd64/$BINARY
GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -o ./bin/linux/arm/$BINARY
GOOS=darwin GOARCH=amd64 go build -trimpath -o ./bin/darwin/amd64/$BINARY
GOOS=windows GOARCH=amd64 go build -trimpath -o ./bin/windows/amd64/$BINARY.exe
md5 ./bin/linux/amd64/$BINARY > ./bin/linux/amd64/$BINARY.md5
md5 ./bin/linux/arm/$BINARY > ./bin/linux/arm/$BINARY.md5
md5 ./bin/darwin/amd64/$BINARY > ./bin/darwin/amd64/$BINARY.md5
md5 ./bin/windows/amd64/$BINARY > ./bin/windows/amd64/$BINARY.md5
