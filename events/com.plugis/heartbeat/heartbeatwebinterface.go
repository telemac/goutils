package heartbeat

import (
	"bufio"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/webserver"
	"github.com/valyala/fasthttp"
	"time"
)

// HeartbeatWebInterface exposes com.plugis.heartbeat.Sent events
// that are saved in the database
type HeartbeatWebInterface struct {
	natsservice.NatsService
	db       Database
	dbConfig DatabaseConfig
}

func NewHeartbeatWebInterface(dbConfig DatabaseConfig) *HeartbeatWebInterface {
	return &HeartbeatWebInterface{dbConfig: dbConfig}
}

func (svc *HeartbeatWebInterface) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()
	log.Debug("heartbeat-web-interface service started")
	defer log.Debug("heartbeat-web-interface service ended")

	err := svc.db.Open(svc.dbConfig)
	if err != nil {
		log.WithError(err).Error("connect to database")
		return err
	}

	server := webserver.NewFiberServer("/views/", "/static/", 8080)

	server.AddTemplateDataProvider("heartbeats", func(c *fiber.Ctx) (fiber.Map, error) {
		records, err := svc.db.getHeartbeats(true)
		return fiber.Map{
			"heartbeats": records,
		}, err
	})

	server.App.Get("/sse/events", func(c *fiber.Ctx) error {
		ctx := c.Context()
		ctx.SetContentType("text/event-stream")
		ctx.Response.Header.Set("Cache-Control", "no-cache")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.Set("Transfer-Encoding", "chunked")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

		ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			log.WithField("ip", ctx.RemoteAddr()).Info("server sent event connected")
			defer log.WithField("ip", ctx.RemoteAddr().String()).Info("server sent event disconnected")

			var i int
			for {
				i++
				msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
				_, err := fmt.Fprintf(w, "event: message\ndata: Message: %s\n\n", msg)
				if err != nil {
					log.WithError(err).Warn("Fprintf in /sse/events")
					return
				}
				fmt.Println(msg)
				err = w.Flush()
				if err != nil {
					log.WithError(err).Warn("Flush in /sse/events")
					return
				}
				time.Sleep(1 * time.Second)
			}
		}))

		return nil
	})

	err = server.Run(ctx)
	if err != nil {
		log.WithError(err).Error("start web server")
	}

	return err
}
