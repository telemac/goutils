package mqtt

import (
	"context"
	"testing"
	"time"

	"github.com/eclipse/paho.golang/paho"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	assert := assert.New(t)
	logrus.SetLevel(logrus.TraceLevel)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	mqttClient := NewMqttClient(MqttParams{
		ServerURL: "tcp://colorbeam:1883",
		ClientID:  "mqtttest",
		//InitialSubscriptions: &paho.Subscribe{
		//	Subscriptions: map[string]paho.SubscribeOptions{
		//		"colorbeam/load/store":        {QoS: 1},
		//		"colorbeam/load/+/status":     {QoS: 1},
		//		"colorbeam/load/+/transition": {QoS: 1},
		//		"colorbeam/+/+/heartbeat":     {QoS: 1},
		//		"colorbeam/building":          {QoS: 1},
		//		"#":                           {QoS: 1},
		//	},
		//},
	})
	err := mqttClient.Connect(ctx)
	assert.NoError(err)
	logrus.WithError(err).Info("connected")

	err = mqttClient.Subscribe(ctx, &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			"colorbeam/load/store": {QoS: 1},
		},
	})

	cancelled := false
	for !cancelled {
		select {
		case msg := <-mqttClient.IncomingMessages():
			logrus.WithField("message", msg).Debug("mqtt service received message")
			if msg.Topic == "colorbeam/load/store" {
				err = mqttClient.Subscribe(ctx, &paho.Subscribe{
					Subscriptions: map[string]paho.SubscribeOptions{
						"colorbeam/load/+/status":     {QoS: 1},
						"colorbeam/load/+/transition": {QoS: 1},
					},
				})
			}
		case <-mqttClient.Done():
			cancelled = true
		}
	}

	// wait
	mqttClient.Done()
	logrus.Info("mqttClient.Done")

	ctx, _ = context.WithTimeout(context.TODO(), time.Second*3)
	mqttClient.Disconnect(ctx)
	mqttClient.Close()

	logrus.WithError(err).Info("closed")

}
