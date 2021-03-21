package main

import (
	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
	"github.com/telemac/goutils/task"
	"github.com/telemac/goutils/webserver"
	"time"
)

func main() {
	ctx, cancel := task.NewCancellableContext(time.Second * 15)
	defer cancel()

	log.SetLevel(log.TraceLevel)

	server := webserver.NewFiberServer("/views/", "/static/", 3000)

	server.AddTemplateDataProvider("colors", func(c *fiber.Ctx) (fiber.Map, error) {
		return fiber.Map{
			"Title":  "AddTemplateDataProvider color",
			"Colors": []string{"red", "green", "blue", "orange", "cyan", "magenta"},
		}, nil
	})

	err := server.Run(ctx)
	if err != nil {
		log.WithError(err).Fatal("start web server")
	}

	//err := webserver.RunFiberServer(ctx, "/views/", "/static/", 3000)
}
