package heartbeat

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"github.com/telemac/goutils/stacktrace"
	"sync"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/gofiber/fiber/v2"
	"github.com/telemac/goutils/natsevents"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/webserver"
	"github.com/valyala/fasthttp"
)

// HeartbeatWebInterface exposes com.plugis.heartbeat.Sent events
// that are saved in the database
type HeartbeatWebInterface struct {
	natsservice.NatsService
	db               Database
	mysqlConfig      natsservice.MysqlConfig
	sseChannels      map[*fiber.Ctx]chan Sent
	sseChannelsMutex sync.RWMutex
}

func NewHeartbeatWebInterface(mysqlConfig natsservice.MysqlConfig) *HeartbeatWebInterface {
	return &HeartbeatWebInterface{mysqlConfig: mysqlConfig,
		sseChannels: make(map[*fiber.Ctx]chan Sent),
	}
}

// OnHeartbeatSent is called each time a heartbeat event is received, must return as fast as possible
func (svc *HeartbeatWebInterface) OnHeartbeatSent(heartbeatSent Sent) {
	log := svc.Logger()
	svc.sseChannelsMutex.RLock()
	defer svc.sseChannelsMutex.RUnlock()
	for _, sentChannel := range svc.sseChannels {
		select {
		case sentChannel <- heartbeatSent:
			log.WithField("heartbeat", heartbeatSent).Trace("enquele heartbeat")
		default:
			log.WithField("stack", stacktrace.GetStackTrace()).Warn("SSE channel is full")
		}
	}
}

func (svc *HeartbeatWebInterface) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()
	log.Debug("heartbeat-web-interface service started")
	defer log.Debug("heartbeat-web-interface service ended")

	err := svc.db.Open(svc.mysqlConfig)
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

		log := svc.Logger().WithField("endpoint", "/sse/events")

		ctx.SetContentType("text/event-stream")
		ctx.Response.Header.Set("Cache-Control", "no-cache")
		ctx.Response.Header.Set("Connection", "keep-alive")
		ctx.Response.Header.Set("Transfer-Encoding", "chunked")
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

		ctx.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			remoteAddr := ctx.RemoteAddr().String()
			log.WithField("ip", remoteAddr).Info("server sent event connected")
			// register a channel to receive HeartbeatSent events
			svc.sseChannelsMutex.Lock()
			svc.sseChannels[c] = make(chan Sent, 10)
			svc.sseChannelsMutex.Unlock()
			defer func() {
				svc.sseChannelsMutex.Lock()
				delete(svc.sseChannels, c)
				svc.sseChannelsMutex.Unlock()
				log.WithField("ip", remoteAddr).Info("server sent event disconnected")
			}()

			for {
				select {
				case <-ctx.Done(): // disconnect
					return
				case event := <-svc.sseChannels[c]:
					sseData, err := json.Marshal(event)
					if err != nil {
						log.WithError(err).Warn("decode heartbeat sent event")
						return
					}

					//_, err = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", "message", sseData)
					_, err = w.WriteString("data: " + string(sseData) + "\n\n")
					if err != nil {
						log.WithError(err).Warn("Write event")
						return
					}
					err = w.Flush()
					if err != nil {
						log.WithError(err).Warn("Flush")
						return
					}
				}
			}
		}))

		return nil
	})

	/* sample curl request
	curl -X POST --location "http://localhost:8080/cloudevents/send" -H "Content-Type: application/json" -d "{
	          \"topic\": \"com.plugis.browser\",
	          \"request\": true,
	          \"timeout\": 5,
	          \"data\": {\"url\": \"https://www.youtube.com\"},
	          \"id\": \"123\",
	          \"source\": \"web\",
	          \"specversion\": \"1.0\",
	          \"type\": \"com.plugis.browser.open\"
	        }"
	*/
	server.App.Post("/cloudevents/send", func(c *fiber.Ctx) error {
		ctx := c.Context()

		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

		sendResult := func(obj interface{}, err error) {
			type Result struct {
				Data  interface{} `json:"data"`
				Error error       `json:"error"`
			}
			type Response struct {
				Result `json:"result"`
			}
			response := Result{obj, err}
			if err != nil {
				c.SendStatus(500)
			} else {
				c.SendStatus(200)
			}
			c.JSON(response)
		}

		var ce cloudevents.Event

		// use own structure for event to allow missing specversion
		//body := c.Body()
		/*
			type SendCloudEvent struct {
				Topic           string      `json:"topic"`
				Type            string      `json:"type"`
				Request         bool        `json:"request"`
				Timeout         int         `json:"timeout"`
				Source          string      `json:"source"`
				Specversion     string      `json:"specversion,omitempty"`
				Datacontenttype string      `json:"datacontenttype"`
				Data            interface{} `json:"data"`
			}
			var sendCloudEvent SendCloudEvent
		*/
		c.Request().Header.Set("Content-Type", "application/json")
		err := c.BodyParser(&ce)

		if err != nil {
			log.WithError(err).Warn("parse cloudEvent from body")
			sendResult(nil, err)
			return err
		}
		// fill missing cloudEvent fields
		natsevents.EventFillDefaults(&ce)

		err = ce.Validate()
		if err != nil {
			log.WithError(err).Warn("validate cloudEvent")
			sendResult(nil, err)
			return err
		}

		extensions := ce.Extensions()
		topic, err := types.ToString(extensions["topic"])
		if err != nil {
			log.WithError(err).Warn("get cloudEvent topic")
			sendResult(nil, err)
			return err
		}

		request := false
		val, ok := extensions["request"]
		if ok {
			request, err = types.ToBool(val)
			if err != nil {
				log.WithError(err).Warn("get cloudEvent request")
				sendResult(nil, err)
				return err
			}
		}

		var timeout int32 = 60 // default 60sec timeout
		val, ok = extensions["timeout"]
		if ok {
			timeout, err = types.ToInteger(val)
			if err != nil {
				log.WithError(err).Warn("get cloudEvent timeout")
				sendResult(nil, err)
				return err
			}
		}

		duration := time.Second * time.Duration(timeout)
		// call cancel function when done
		ctxTimeout, cancel := context.WithTimeout(context.TODO(), duration)
		defer cancel()
		if request {
			returnedEvent, err := svc.Transport().Request(ctxTimeout, &ce, topic, duration)
			if err != nil {
				log.WithError(err).Warn("cloudEvent request")
				sendResult(nil, err)
				return err
			}
			sendResult(returnedEvent, nil)
		} else {
			err = svc.Transport().Send(ctx, &ce, topic)
			if err != nil {
				log.WithError(err).Warn("send cloudEvent")
				sendResult(nil, err)
			}
			sendResult(struct{}{}, nil)
		}

		log.Tracef("cloudEvent = %+v\n", ce)
		log.Tracef("topic = %s\n", topic)
		log.Tracef("request = %v\n", request)
		log.Tracef("timeout = %d\n", timeout)

		return nil
	})

	server.App.Get("/heartbeats/comments", func(c *fiber.Ctx) error {
		ctx := c.Context()
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Cache-Control")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		// TODO : return mac with non empty comments
		return errors.New("/heartbeats/comments not implemented ")
	})

	err = server.Run(ctx)
	if err != nil {
		log.WithError(err).Error("start web server")
	}

	return err
}
