# event-sender

## sample
```shell

# open browser
./event-sender -type 'com.plugis.browser.open' -data '{"url":"https://google.fr"}' -server 'https://nats1.plugis.com'

# shell command
./event-sender -type 'com.plugis.shell.command' -data '{"command": ["df","-lh","/"]}' -topic "com.plugis.shell.B8:27:EB:E7:B3:4E" -request
./event-sender -type 'com.plugis.shell.command' -data '{"command": ["sudo","apt","update"]}' -topic "com.plugis.shell.B8:27:EB:E7:B3:4E" -request
./event-sender -type 'com.plugis.shell.command' -data '{"command": ["sudo","apt","-y","upgrade"]}' -topic "com.plugis.shell" -request
./event-sender -type 'com.plugis.shell.command' -data '{"command": ["hostnamectl"]}' -topic "com.plugis.shell" -request
./event-sender -type 'com.plugis.shell.command' -data '{"command": ["sudo","arp-scan","--interface=eth0","--localnet"]}' -topic "com.plugis.shell" -request


```
