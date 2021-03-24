package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/events/com.plugis/heartbeat"
	"github.com/telemac/goutils/events/com.plugis/service"
	"github.com/telemac/goutils/events/com.plugis/shell"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/task"
	"github.com/telemac/goutils/updater"
	"time"
)

type CommandLineParams struct {
	Install bool
}

var commandLineParams CommandLineParams

func main() {
	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	// get command line params
	flag.BoolVar(&commandLineParams.Install, "install", false, "install as service")
	flag.Parse()

	logrus.SetLevel(logrus.TraceLevel)

	// install as service and exit if -install parameter present
	selfInstallService := &service.SelfInstallService{
		ServiceName: "remote-access",
	}
	if commandLineParams.Install {
		err := selfInstallService.Install()
		if err != nil {
			logrus.WithError(err).Error("installint service")
		}
		return
	}

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
	}

	//servicesRepository, err := natsservice.NewNatsServiceRepository("remote-access", "nats://cloud1.idronebox.com:443", "trace")
	servicesRepository, err := natsservice.NewNatsServiceRepository("remote-access", "nats://nats1.plugis.com:443", "trace")
	if err != nil {
		logrus.WithError(err).Fatal("create nats service repository")
	}
	defer servicesRepository.Close(time.Second * 10)

	servicesRepository.Logger().Info("remote-access service starting")

	// auto install service
	servicesRepository.Start(ctx, selfInstallService)

	// start heartbeat service
	servicesRepository.Start(ctx, &heartbeat.HeartbeatSender{
		Period:       55,
		RandomPeriod: 4,
	})

	// start shell service
	servicesRepository.Start(ctx, &shell.ShellService{})

	// browser service
	//servicesRepository.Start(ctx, &browser.BrowserService{})

	servicesRepository.WaitUntilAllDone()

	servicesRepository.Logger().Info("remote-access service ending")
}
