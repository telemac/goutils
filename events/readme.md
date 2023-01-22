# Cloud Events

## send manually with nats-pub

send a minimal event manually with nats-pub

```shell
nats pub -s 'https://nats1.plugis.com' 'com.plugis.browser' '{"type": "com.plugis.browser.open","data": {"url": "https://www.plugis.com"}, "id": "123","source": "manual","specversion": "1.0"}'
nats pub -s 'https://nats1.plugis.com' 'com.plugis.browser' '{"type": "com.plugis.browser.open","data": {"url": "https://www.google.fr"}, "id": "123","source": "manual","specversion": "1.0"}'
```

## use new nats client
[natscli](https://github.com/nats-io/natscli)

```shell
# create local context
nats context add local --server localhost:4222 --description "Localhost" --select
# add user/password
export EDITOR=vim
nats context edit local

# use local context
nats context select local

# test round-trip time
nats rtt

# get jetstream account info
nats account info

# see jetstream events
nats event --js-advisory

# get event json schema
nats schema show io.nats.jetstream.api.v1.stream_create_response


# subscribe on all topics
nats sub '>'

```

## MQTT client with nats
mosquitto_sub -h update.idronebox.com -t '#' -v -d -u <user> -P <password>

## AsyncAPI with CloudEvents
[Simulating CloudEvents with AsyncAPI and Microcks](https://developers.redhat.com/articles/2021/06/02/simulating-cloudevents-asyncapi-and-microcks#)
