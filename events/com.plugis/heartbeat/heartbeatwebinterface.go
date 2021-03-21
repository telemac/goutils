package heartbeat

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/webserver"
)

// HeartbeatWebInterface exposes com.plugis.heartbeat.Sent events
// that are saved in the database
type HeartbeatWebInterface struct {
	natsservice.NatsService
	db Database
}

func (svc *HeartbeatWebInterface) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()
	log.Debug("heartbeat-web-interface service started")
	defer log.Debug("heartbeat-web-interface service ended")

	// premare database
	dbConfig := DatabaseConfig{
		DBHost: "127.0.0.1",
		DBname: "plugis",
		DBuser: "root",
		DBpass: "telemac",
		DBPort: 3306,
	}
	err := svc.db.Open(dbConfig)
	if err != nil {
		log.WithError(err).Error("connect to database")
		return err
	}

	server := webserver.NewFiberServer("/views/", "/static/", 8080)

	server.AddTemplateDataProvider("heartbeats", func(c *fiber.Ctx) (fiber.Map, error) {
		records, err := svc.db.getHeartbeat(true)
		return fiber.Map{
			"heartbeats": records,
		}, err
	})

	err = server.Run(ctx)
	if err != nil {
		log.WithError(err).Error("start web server")
	}

	return err
}
