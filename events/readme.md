# Cloud Events

## send manually with nats-pub

send a minimal event manually with nats-pub

```shell
nats-pub -s 'https://nats1.plugis.com' 'com.plugis.browser' '{"type": "com.plugis.browser.open","data": {"url": "https://www.plugis.com"}, "id": "123","source": "manual","specversion": "1.0"}'
nats-pub -s 'https://nats1.plugis.com' 'com.plugis.browser' '{"type": "com.plugis.browser.open","data": {"url": "https://www.google.fr"}, "id": "123","source": "manual","specversion": "1.0"}'
```
