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

func (svc *SelfInstallService) Install(arguments []string) error {
	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        svc.ServiceName,
		Arguments:   arguments,
		DisplayName: svc.ServiceName,
		Description: "remote-access service.",
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"},
		Option: options,
	}

	s, err := service.New(svc, svcConfig)
	if err != nil {
		return err
	}

	err = s.Install()
	if err == nil {
		// start the installed service
		err = s.Start()
		if err != nil {
			return err
		}

	}
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		return err
	}
	return nil
}

func (svc *SelfInstallService) Uninstall() error {
	svcConfig := &service.Config{
		Name:        svc.ServiceName,
		DisplayName: svc.ServiceName,
	}

	s, err := service.New(svc, svcConfig)
	if err != nil {
		return err
	}

	err = s.Uninstall()
	if err == nil {
		return err
	}
	return s.Stop()
}

func (svc *SelfInstallService) Run(ctx context.Context, params ...interface{}) error {
	svc.Logger().Debug("serSelfInstallServicevice started")
	defer svc.Logger().Debug("SelfInstallService ended")

	interactive := service.Interactive()

	if interactive {
		err := svc.Install(nil)
		if err != nil {
			return err
		}
	}

	//err = s.Run()
	//if err != nil {
	//	log.WithError(err).Error("run service")
	//}

	return nil
}
