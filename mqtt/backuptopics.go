package mqtt

import (
	"context"
	"fmt"

	"github.com/eclipse/paho.golang/paho"
	"github.com/sirupsen/logrus"
)

// TopicsBackup holds the mqtt topics values
type TopicsBackup map[string]string

// BackupTopics connects to an mqtt server and writes the topics and their content into the writer
func BackupTopics(ctx context.Context, server string, topics []string) (TopicsBackup, error) {

	topicsBackup := make(TopicsBackup)

	var subscriptions = make(map[string]paho.SubscribeOptions)
	for _, topic := range topics {
		subscriptions[topic] = paho.SubscribeOptions{
			QoS: 1,
		}
	}

	mqttClient := NewMqttClient(MqttParams{
		ServerURL: server,
		ClientID:  "mqtt.BackupTopics",
		InitialSubscriptions: &paho.Subscribe{
			Subscriptions: subscriptions,
		},
	})
	err := mqttClient.Connect(ctx)
	if err != nil {
		return topicsBackup, fmt.Errorf("connect to mqtt server %s : %w", server, err)
	}

	defer func() {
		mqttClient.Disconnect(ctx)
		mqttClient.Close()
	}()

	cancelled := false
	for !cancelled {
		select {
		case msg := <-mqttClient.IncomingMessages():
			topicsBackup[msg.Topic] = string(msg.Payload)
			logrus.WithField("message", msg).Info("mqtt service received message")
		case <-mqttClient.Done():
			cancelled = true
		}
	}

	return topicsBackup, nil
}
