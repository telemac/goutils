GOOS=linux GOARCH=amd64 go build -trimpath -o ./bin/remote-access_linux_amd64
GOOS=linux GOARCH=arm GOARM=7 go build -trimpath -o ./bin/remote-access_linux_arm
GOOS=darwin GOARCH=amd64 go build -trimpath -o ./bin/remote-access_darwin_amd64
GOOS=windows GOARCH=amd64 go build -trimpath -o ./bin/remote-access_windows.exe

rsync -avP ./bin/remote-access_linux_amd64 cloud.plugis.com:
rsync -avP ./bin/remote-access_linux_amd64 colorbeam@cbna:
rsync -avP ./bin/remote-access_linux_arm colorbeam@colorbeam:
