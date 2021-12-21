# tcp_get_command

the tcp_get_command service listens on a tcp port and publishes the received data in a cloud event on the topic com.plugis.tcp-get-command.message

```shell
# in one terminal
cd github.com/telemac/goutils/events/com.plugis/tcpgetcommand/cmd/tcp_get_command_sample
go run main.go -log trace

# in another
echo -n "000042<MESSAGE_TYPE>=WATCHDOG<MODULE>=WATCH_DC09" | nc localhost 10

# in a third
 mosquitto_sub -h cloud1.idronebox.com -t 'com/#' -v

```
