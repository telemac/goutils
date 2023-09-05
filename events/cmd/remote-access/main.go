package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/ansible"
	"github.com/telemac/goutils/events/com.plugis/browser"
	"github.com/telemac/goutils/events/com.plugis/heartbeat"
	"github.com/telemac/goutils/events/com.plugis/service"
	"github.com/telemac/goutils/events/com.plugis/shell"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/net"
	"github.com/telemac/goutils/remote-access"
	"github.com/telemac/goutils/task"
	"github.com/telemac/goutils/updater"
	"os"
	"strings"
	"time"
)

func main() {
	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	// get command line params
	var commandLineParams remote_access.CommandLineParams
	flag.BoolVar(&commandLineParams.Install, "install", false, "install as service")
	flag.BoolVar(&commandLineParams.Uninstall, "uninstall", false, "uninstall the service")
	flag.BoolVar(&commandLineParams.Update, "update", true, "self update at startup")
	flag.StringVar(&commandLineParams.Log, "log", "warn", "log level (trace|debug|info|warn|error)")
	flag.StringVar(&commandLineParams.NatsServers, "nats", "", "nats server urls separated by ,")

	flag.Parse()

	logrus.SetLevel(logrus.TraceLevel)

	// install as service and exit if -install parameter present
	selfInstallService := &service.SelfInstallService{
		ServiceName: "remote-access",
	}
	if commandLineParams.Install {
		var arguments []string
		if commandLineParams.NatsServers != "" {
			arguments = append(arguments, "-nats", commandLineParams.NatsServers)
		}

		err := selfInstallService.Install(arguments)
		if err != nil {
			logrus.WithError(err).Error("installing service")
		} else {
			logrus.Info("service installed")
		}
		return
	}

	if commandLineParams.Uninstall {

		err := selfInstallService.Uninstall()
		if err != nil {
			logrus.WithError(err).Error("uninstalling service")
		} else {
			logrus.Info("service uninstalled")
		}
		return
	}

	if commandLineParams.Update {
		// self update binary
		logrus.Info("check for update")
		selfUpdater, err := updater.NewSelfUpdater("https://update.plugis.com/", "")
		if err != nil {
			logrus.WithError(err).Error("self update creation")
		}
		needsUpdate, err := selfUpdater.NeedsUpdate()
		if err != nil {
			logrus.WithError(err).Error("check if update needed")
		}
		if needsUpdate {
			logrus.Info("start self updating...")
		} else if err == nil {
			logrus.Info("binary is up-to-date")
		}

		updated, err := selfUpdater.SelfUpdate(false)

		if err != nil {
			logrus.WithError(err).Error("self update")
		}
		if updated {
			logrus.Info("is updated, must restart")
			os.Exit(1)
		}
	}

	//servicesRepository, err := natsservice.NewNatsServiceRepository("remote-access", "nats://cloud1.idronebox.com:443", "trace")

	var servers []string
	if commandLineParams.NatsServers != "" {
		servers = strings.Split(commandLineParams.NatsServers, ",")
	}

	servers = append(servers, "wss://remote-access:@nats.plugis.cloud:443", "ws://remote-access:@nats.plugis.cloud:8222")
	serversList := strings.Join(servers, ",")
	servicesRepository, err := natsservice.NewNatsServiceRepository("remote-access", serversList, commandLineParams.Log)
	if err != nil {
		logrus.WithError(err).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	macAddress, err := net.GetMACAddress()
	if err != nil {
		servicesRepository.Logger().WithError(err).Error("get mac address")
	}

	servicesRepository.Logger().WithFields(logrus.Fields{"mac": macAddress}).Info("remote-access service starting")

	// auto install service
	servicesRepository.Start(ctx, selfInstallService)

	// start heartbeat service
	servicesRepository.Start(ctx, &heartbeat.HeartbeatSender{
		Period:       55,
		RandomPeriod: 4,
		Meta: map[string]interface{}{
			"version": "0.0.2",
		},
	})

	// com.plugis.shell service
	servicesRepository.Start(ctx, &shell.ShellService{})

	// com.plugis.browser service
	servicesRepository.Start(ctx, &browser.BrowserService{})

	// com.plugis.ansible service
	servicesRepository.Start(ctx, &ansible.AnsibleService{})

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info("remote-access service ending")
}
