package service

import (
	"context"
	"strings"

	"github.com/kardianos/service"
	"github.com/telemac/goutils/natsservice"
)

// SelfInstallService installs the current executable as a service on linux
type SelfInstallService struct {
	natsservice.NatsService
	ServiceName string
}

func (svc *SelfInstallService) Start(s service.Service) error {
	svc.Logger().Debug("start called")
	return nil
}

func (svc *SelfInstallService) Stop(s service.Service) error {
	svc.Logger().Debug("stop called")
	return nil
}

func (svc *SelfInstallService) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()
	log.Debug("serSelfInstallServicevice started")
	defer log.Debug("SelfInstallService ended")

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        svc.ServiceName,
		DisplayName: svc.ServiceName,
		Description: "remote-access service.",
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"},
		Option: options,
	}

	interactive := service.Interactive()
	s, err := service.New(svc, svcConfig)
	if err != nil {
		log.Fatal("install service")
	}
	if interactive {
		log.Debug("install service")
		err = s.Install()
		if err == nil {
			// start the installed service
			err = s.Start()
			if err != nil {
				log.WithError(err).Error("start service")
			}

		}
		if err != nil && !strings.Contains(err.Error(), "already exists") {
			log.WithError(err).Error("install service")
		}
	}

	//err = s.Run()
	//if err != nil {
	//	log.WithError(err).Error("run service")
	//}

	return nil
}
