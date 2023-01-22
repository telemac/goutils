# remote-access

## service path on MacOS :
/Library/LaunchDaemons/remote-access.plist

## deploy
wget update.plugis.com/linux/arm/remote-access
wget update.plugis.com/linux/amd64/remote-access
wget update.plugis.com/darwin/amd64/remote-access

chmod +x remote-access
sudo ./remote-access -install
sudo service remote-access status

## show remote-access service logs
# sudo journalctl -x -u remote-access
