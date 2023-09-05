package maintenance

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/heartbeat"
	"github.com/telemac/goutils/natsevents"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/net"
	"os/exec"
)

type MaintenanceService struct {
	natsservice.NatsService
	Config           MaintenanceServiceConfig
	heartbeatService *heartbeat.HeartbeatSender
}

// MaintenanceServiceConfig is the structure of the configuration for the MaintenanceService
type MaintenanceServiceConfig struct {
	SSHServer   string `json:"ssh_server"`
	SSHPort     int    `json:"ssh_port"`
	SSHUser     string `json:"ssh_user"`
	SSHPassword string `json:"ssh_password"`
	SSHKeyPath  string `json:"ssh_key_path"`
}

// MaintenanceStartParams is the structure of the event com.plugis.maintenance.start
type MaintenanceStartParams struct {
	Port int `json:"port"`
}

type MaintenanceResponse struct {
	Response string `json:"response,omitempty"`
	Error    error  `json:"error,omitempty"`
}

// NewMaintenanceService creates a new MaintenanceService
func NewMaintenanceService(config MaintenanceServiceConfig, heartbeatService *heartbeat.HeartbeatSender) *MaintenanceService {
	svc := &MaintenanceService{
		Config:           config,
		heartbeatService: heartbeatService,
	}
	return svc
}

func (svc *MaintenanceService) StartMaintenance(startParams *MaintenanceStartParams) error {

	if startParams.Port <= 0 {
		return errors.New("maintenance port not set")
	}

	// "2000:localhost:22"
	NRStr := fmt.Sprintf("%d:localhost:22", startParams.Port)
	// "maintenance@nats.plugis.cloud"
	sshConnectionStr := fmt.Sprintf("%s@%s", svc.Config.SSHUser, svc.Config.SSHServer)
	// "-p2100"
	sshPortStr := fmt.Sprintf("-p%d", svc.Config.SSHPort)

	// sshpass -p maintenance ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no -o ServerAliveInterval=100 -NR 2000:localhost:22 maintenance@cb-na.cloud -p2100
	shellCommand := []string{"sshpass", "-p", svc.Config.SSHPassword, "ssh", "-o", "UserKnownHostsFile=/dev/null", "-o", "StrictHostKeyChecking=no", "-o", "ServerAliveInterval=100", "-NR", NRStr, sshConnectionStr, sshPortStr}

	var cmd *exec.Cmd
	cmd = exec.Command(shellCommand[0], shellCommand[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		svc.Logger().WithError(err).WithFields(logrus.Fields{
			"command": shellCommand,
			"output":  string(out),
		}).Warn("run command")
		return err
	}

	return nil
}

func (svc *MaintenanceService) eventHandler(topic string, receivedEvent *event.Event, payload []byte, err error) (*event.Event, error) {
	// check if no error on cloud event formatting
	if err != nil {
		svc.Logger().WithFields(logrus.Fields{
			"topic":   topic,
			"event":   receivedEvent,
			"payload": string(payload),
			"error":   err,
		}).Error("receive cloud event")
		return nil, err
	}

	switch receivedEvent.Type() {
	case "com.plugis.maintenance.start":
		var params MaintenanceStartParams
		err = receivedEvent.DataAs(&params)
		if err != nil {
			svc.Logger().WithError(err).Warn("decode ShellCommandParams")
			return nil, err
		}

		// TODO : set maintenance meta on heartbeat in a thread save way
		svc.heartbeatService.AddMeta("maintenance", params.Port, true)
		err = svc.StartMaintenance(&params)
		svc.heartbeatService.RemoveMeta("maintenance", true)

		if err != nil {
			svc.Logger().WithError(err).Warn("start maintenance")
			return nil, err
		}
		// TODO : signal maintenance started and save

		responseEvent := natsevents.NewEvent("", "com.plugis.shell.response", params)
		responseEvent.SetExtension("responsefor", receivedEvent.ID())
		responseEvent.SetSource("com.plugis.shell")

		return responseEvent, err

	default:
		svc.Logger().WithFields(logrus.Fields{
			"topic": topic,
			"type":  receivedEvent.Type(),
			"event": receivedEvent,
		}).Warn("unknown event type")
	}

	return nil, errors.New("unattainable code")
}

func (svc MaintenanceService) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()

	log.Debug("MaintenanceService started")
	defer log.Debug("MaintenanceService ended")

	var err error

	// register eventHandler for event reception
	mac, err := net.GetMACAddress()
	if err != nil {
		svc.Logger().WithError(err).Error("get mac accress")
	}
	topic := "com.plugis.maintenance." + mac
	err = svc.Transport().RegisterHandler(svc.eventHandler, topic)
	if err != nil {
		svc.Logger().WithError(err).Error("register event handler")
		return err
	}

	<-ctx.Done()

	return nil
}
