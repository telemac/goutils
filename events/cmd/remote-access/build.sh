BINARY=remote-access

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -o ./bin/linux/amd64/$BINARY
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -o ./bin/linux/arm/$BINARY
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -o ./bin/darwin/amd64/$BINARY
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -o ./bin/windows/amd64/$BINARY.exe
md5 ./bin/linux/amd64/$BINARY > ./bin/linux/amd64/$BINARY.md5
md5 ./bin/linux/arm/$BINARY > ./bin/linux/arm/$BINARY.md5
md5 ./bin/darwin/amd64/$BINARY > ./bin/darwin/amd64/$BINARY.md5
md5 ./bin/windows/amd64/$BINARY.exe > ./bin/windows/amd64/$BINARY.exe.md5
