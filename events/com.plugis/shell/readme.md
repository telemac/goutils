# com.plugis.shell cloud events

Topic : com.plugis.shell

Cloud events :
- Type : com.plugis.shell.command
- parameters : {"command": ["df","-lh","/"]}
- sample : 
nats-pub -s 'https://nats1.plugis.com' 'com.plugis.shell' '{"type": "com.plugis.shell.command","data": {"command": ["df","-lh","/"]}, "id": "123","source": "manual","specversion": "1.0"}'

