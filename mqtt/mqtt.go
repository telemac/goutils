package mqtt

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/logger"
)

type MqttParams struct {
	ServerURL            string
	ClientID             string
	InitialSubscriptions *paho.Subscribe
}

type MqttClient struct {
	log               *logrus.Entry
	params            MqttParams
	connectionManager *autopaho.ConnectionManager
	incomingMessages  chan *paho.Publish
}

func NewMqttClient(params MqttParams) *MqttClient {
	return &MqttClient{
		params: params,
	}
}

// onMessage is called when a message is received
func (mqttClient *MqttClient) onMessage(msg *paho.Publish) {
	mqttClient.log.WithFields(logrus.Fields{
		"topic":   msg.Topic,
		"payload": string(msg.Payload),
		"src":     "MqttClient.onMessage",
	}).Trace("received message")

	// push message to channel
	select {
	case mqttClient.incomingMessages <- msg:
	default:
		mqttClient.log.WithFields(logrus.Fields{
			"topic":   msg.Topic,
			"payload": string(msg.Payload),
		}).Warn("push mqtt message to incoming channel")
	}

}

// logger implements the paho.Logger interface
type pahoLogger struct {
	log *logrus.Entry
}

// Println is the library provided NOOPLogger's
// implementation of the required interface function()
func (l pahoLogger) Println(v ...interface{}) {
	l.log.Trace(v...)
}

// Printf is the library provided NOOPLogger's
// implementation of the required interface function(){}
func (l pahoLogger) Printf(format string, v ...interface{}) {
	l.log.Tracef(format, v...)
}

func (mqttClient *MqttClient) Connect(ctx context.Context) error {
	// build logger
	mqttClient.log = logger.FromContext(ctx, true).WithFields(logrus.Fields{
		"url": mqttClient.params.ServerURL,
	})
	mqttClient.log.Debug("mqtt connect")

	// parse mqtt url
	u, err := url.Parse(mqttClient.params.ServerURL)
	if err != nil {
		return err
	}

	uid, err := uuid.NewV4()
	if err != nil {
		mqttClient.log.Error("generate uuid v4")
	}

	// build autopaho configuration
	cliCfg := autopaho.ClientConfig{
		//Debug:             mqttClient.log.WithField("src", "autoPaho"),
		//PahoDebug:         pahoLogger{log: mqttClient.log.WithField("src", "paho")},
		BrokerUrls:        []*url.URL{u},
		KeepAlive:         10,
		ConnectRetryDelay: time.Second * 10,
		OnConnectError: func(err error) {
			mqttClient.log.WithError(err).Warn("error whilst attempting connection")
		},
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			mqttClient.log.Debug("mqtt connection up")
			if mqttClient.params.InitialSubscriptions != nil {
				_, err = cm.Subscribe(context.Background(), mqttClient.params.InitialSubscriptions)
				if err != nil {
					mqttClient.log.WithError(err).Warn("failed to subscribe")
					return
				}
				mqttClient.log.Info("mqtt subscription ok")
			}
		},
		ClientConfig: paho.ClientConfig{
			ClientID: mqttClient.params.ClientID + "." + uid.String(),
			Router: paho.NewSingleHandlerRouter(func(m *paho.Publish) {
				mqttClient.onMessage(m)
			}),
			OnClientError: func(err error) { fmt.Printf("server requested disconnect: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					fmt.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					fmt.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}

	mqttClient.connectionManager, err = autopaho.NewConnection(ctx, cliCfg)
	if err != nil {
		return err
	}

	mqttClient.incomingMessages = make(chan *paho.Publish, 1000)

	err = mqttClient.connectionManager.AwaitConnection(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (mqttClient *MqttClient) Disconnect(ctx context.Context) error {
	return mqttClient.connectionManager.Disconnect(ctx)
}

func (mqttClient *MqttClient) Done() <-chan struct{} {
	return mqttClient.connectionManager.Done()
}

func (mqttClient *MqttClient) Close() {
	close(mqttClient.incomingMessages)
}

func (mqttClient *MqttClient) IncomingMessages() <-chan *paho.Publish {
	return mqttClient.incomingMessages
}

func (mqttClient *MqttClient) Publish(ctx context.Context, p *paho.Publish) error {
	_, err := mqttClient.connectionManager.Publish(ctx, p)
	return err
}

func (mqttClient *MqttClient) Subscribe(ctx context.Context, s *paho.Subscribe) error {
	_, err := mqttClient.connectionManager.Subscribe(ctx, s)
	return err
}
