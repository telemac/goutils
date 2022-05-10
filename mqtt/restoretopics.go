package mqtt

import (
	"context"
	"fmt"
	"time"

	"github.com/eclipse/paho.golang/paho"
	"github.com/telemac/goutils/task"
)

// RestoreTopics publishes topics from  topicsBackup as retaines
func RestoreTopics(ctx context.Context, server string, topicsBackup TopicsBackup) error {
	mqttClient := NewMqttClient(MqttParams{
		ServerURL: server,
		ClientID:  "mqtt.RestoreTopics",
	})
	err := mqttClient.Connect(ctx)
	if err != nil {
		return fmt.Errorf("connect to mqtt server %s : %w", server, err)
	}

	defer func() {
		mqttClient.Disconnect(ctx)
		mqttClient.Close()
	}()

	for topic, payload := range topicsBackup {
		err = mqttClient.Publish(ctx, &paho.Publish{
			Topic:   topic,
			Payload: []byte(payload),
			Retain:  true,
			QoS:     1,
		})
		if err != nil {
			return err
		}
	}

	task.Sleep(ctx, time.Millisecond*100)
	return nil
}
