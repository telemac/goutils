package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/ansibleutils"
	"github.com/telemac/goutils/logger"
	"github.com/telemac/goutils/task"
	"strings"
	"time"
)

type CommandLineParams struct {
	Playbook string
	Log      string
	ansibleutils.AnsibleParams
}

func main() {
	ctx, cancel := task.NewCancellableContext(time.Second * 5)
	defer cancel()

	var commandLineParams CommandLineParams
	flag.StringVar(&commandLineParams.Playbook, "playbooks", "/tmp/playbook.yml", "playbook files separated by a space")
	flag.StringVar(&commandLineParams.Log, "log", "info", "log level")
	flag.StringVar(&commandLineParams.Inventory, "inventory", "/tmp/inventory.yml", "inventory file")
	flag.StringVar(&commandLineParams.Roles, "roles", "", "roles separated by a space")
	flag.StringVar(&commandLineParams.Packages, "packages", "", "packages separated by a space")
	flag.StringVar(&commandLineParams.Base, "base", "", "base directory url for playbooks and inventory")
	flag.Parse()

	// create a logger with the give log level
	log := logger.New(commandLineParams.Log, logrus.Fields{"event-type": "com.plugis.ansible"})
	// add the logger to the main context
	ctx = logger.WithLogger(ctx, log)

	log.Info("ansible playbook runner started")
	defer log.Info("ansible playbook runner finished")

	// split playbook
	commandLineParams.Playbooks = strings.Split(commandLineParams.Playbook, " ")

	a := ansibleutils.New()
	result, err := a.RunPlaybooks(ctx, commandLineParams.AnsibleParams)
	_ = result
	if err != nil {
		log.WithError(err).Error("run playbook")
	}
	//fmt.Println(result)
}
