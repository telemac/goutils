# com.plugis.shell cloud events

Topic : com.plugis.shell

Cloud events :
- Type : com.plugis.shell.command.B8:27:EB:E7:B3:4E
- parameters : {"command": ["df","-lh","/"]}
- sample : 
nats-pub -s 'https://nats1.plugis.com' 'com.plugis.shell.B8:27:EB:E7:B3:4E' '{"type": "com.plugis.shell.command","data": {"command": ["df","-lh","/"]}, "id": "123","source": "manual","specversion": "1.0"}'

event-sender -type 'com.plugis.shell.command' -data '{"command": ["hostname"]}' -topic "com.plugis.shell.B8:27:EB:E7:B3:4E" -request
