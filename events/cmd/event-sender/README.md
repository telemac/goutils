# event-sender

## sample
```shell

# open browser
./event-sender -type 'com.plugis.browser.open' -data '{"url":"https://google.fr"}' -router 'https://nats1.plugis.com'

# shell command
./event-sender -type 'com.plugis.shell.command' -data '{"command": ["df","-lh","/"]}' -topic "com.plugis.shell.B8:27:EB:E7:B3:4E" -request
./event-sender -type 'com.plugis.shell.command' -data '{"command": ["sudo","apt","update"]}' -topic "com.plugis.shell.B8:27:EB:E7:B3:4E" -request
./event-sender -type 'com.plugis.shell.command' -data '{"command": ["sudo","apt","-y","upgrade"]}' -topic "com.plugis.shell.B8:27:EB:E7:B3:4E" -request
./event-sender -type 'com.plugis.shell.command' -data '{"command": ["hostnamectl"]}' -topic "com.plugis.shell.B8:27:EB:E7:B3:4E" -request -timeout 10
./event-sender -type 'com.plugis.shell.command' -data '{"command": ["sudo","arp-scan","--interface=eth0","--localnet"]}' -topic "com.plugis.shell.B8:27:EB:E7:B3:4E" -request

./event-sender -type 'com.plugis.shell.command' -data '{"command": ["sudo","reboot"]}' -topic "com.plugis.shell.B8:27:EB:2C:83:55" -request

./event-sender -type 'com.plugis.shell.command' -data '{"command": ["sh","-c","apt install -y ansible"]}' -topic "com.plugis.shell.8E:82:45:0E:A4:6F" -request -timeout 300


```

## add the command on osX
```shell
sudo ln -s event-sender /usr/local/bin/event-sender
```
