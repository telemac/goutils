# com.plugis.maintenance cloud events

Topic : com.plugis.maintenance.B8:27:EB:E7:B3:4E

Cloud events :
- Type : com.plugis.maintenance.start
- parameters : {"port": 2000}
- sample : 
nats pub -s 'https://nats1.plugis.com' 'com.plugis.maintenance.B8:27:EB:E7:B3:4E' '{"type": "com.plugis.maintenance.start","data": {"port": 2000}, "id": "123","source": "manual","specversion": "1.0"}'

event-sender -type 'com.plugis.maintenance.start' -data '{"port": 2000}' -topic "com.plugis.maintenance.B8:27:EB:E7:B3:4E" -request

event-sender -type 'com.plugis.maintenance.stop' -request

## start maintenance
event-sender -server "nats://nats1.plugis.com:443" -type 'com.plugis.maintenance.start' -data '{"port": 2000}' -topic "com.plugis.maintenance.B8:27:EB:E7:B3:4E" -timeout 10 -request
