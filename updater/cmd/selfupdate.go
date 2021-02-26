package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/telemac/goutils/updater"
	"os"
)

func main() {
	selfUpdater, err := updater.NewSelfUpdater("https://update.plugis.com/", "")
	if err != nil {
		log.WithError(err).Error("self update creation")
	}
	updated, err := selfUpdater.SelfUpdate(false)

	if err != nil {
		log.WithError(err).Error("self update")
	}
	if updated {
		fmt.Println("is updated, must restart")
	}
	fmt.Println("version 0.0.13")
	os.Exit(0)
}
