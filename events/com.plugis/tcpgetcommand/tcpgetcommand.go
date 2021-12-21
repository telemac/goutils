package tcpgetcommand

import (
	"context"
	"strconv"
	"time"

	"github.com/Allenxuxu/gev"
	"github.com/Allenxuxu/gev/connection"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/natsevents"
	"github.com/telemac/goutils/natsservice"
)

type TcpGetCommandConfig struct {
	ListenPort    int    `json:"listen_port"`
	ListenAddress string `json:"listen_address"`
}

type TcpGetCommandService struct {
	natsservice.NatsService
	config TcpGetCommandConfig
}

type TcpGetCommandMessage struct {
	Data       []byte
	RemoteAddr string
}

func NewTcpGetCommandService(config TcpGetCommandConfig) *TcpGetCommandService {
	return &TcpGetCommandService{config: config}
}

func (svc *TcpGetCommandService) OnConnect(c *connection.Connection) {
	svc.Logger().WithFields(logrus.Fields{
		"peer_addr":          c.PeerAddr(),
		"read_buffer_length": c.ReadBufferLength(),
	}).Debug("OnConnect")
	//c.Set("name", "Alexandre")
}

func (svc *TcpGetCommandService) OnMessage(c *connection.Connection, ctx interface{}, data []byte) (out interface{}) {
	// TODO : decode watchcomx message 000042<MESSAGE_TYPE>=WATCHDOG<MODULE>=WATCH_DC09
	svc.Logger().WithFields(logrus.Fields{
		"data":               string(data),
		"peer_addr":          c.PeerAddr(),
		"read_buffer_length": c.ReadBufferLength(),
	}).Debug("OnMessage")

	eventData := TcpGetCommandMessage{
		Data:       data,
		RemoteAddr: c.PeerAddr(),
	}
	heartbeatEvent := natsevents.NewEvent("com.plugis.", "tcp-get-command.message", eventData)
	topic := heartbeatEvent.Type()
	err := svc.Transport().Send(context.TODO(), heartbeatEvent, topic)
	if err != nil {
		svc.Logger().WithFields(logrus.Fields{
			"error":     err,
			"type":      heartbeatEvent.Type(),
			"data":      string(data),
			"peer_addr": c.PeerAddr(),
		}).Error("send event")
	}

	name, ok := c.Get("name")
	if ok {
		svc.Logger().Println("name", name.(string))
	}
	// out = []byte("OK\r\n")
	return
}

func (svc *TcpGetCommandService) OnClose(c *connection.Connection) {
	c.UserBuffer()
	svc.Logger().WithFields(logrus.Fields{
		"peer_addr": c.PeerAddr(),
	}).Debug("OnClose")
}

func (svc *TcpGetCommandService) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()

	log.Debug("tcp-get-command service started")
	defer log.Debug("tcp-get-command service ended")

	s, err := gev.NewServer(svc,
		gev.Address(svc.config.ListenAddress+":"+strconv.Itoa(svc.config.ListenPort)),
		gev.NumLoops(0),
		gev.IdleTime(65*time.Second))

	if err != nil {
		panic(err)
	}

	go s.Start()
	defer s.Stop()

	<-ctx.Done()
	s.Stop()

	return nil
}
