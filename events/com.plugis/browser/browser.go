package browser

import (
	"context"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/natsservice"
)

type BrowserService struct {
	natsservice.NatsService
}

func (svc *BrowserService) eventHandler(topic string, event *event.Event, payload []byte, err error) (*event.Event, error) {
	// check if no error on cloud event formatting
	if err != nil {
		svc.Logger().WithFields(logrus.Fields{
			"topic":   topic,
			"event":   event,
			"payload": string(payload),
			"error":   err,
		}).Error("receive cloud event")
		return nil, err
	}

	switch event.Type() {
	case "com.plugis.browser.open":
		type OpenBrowserParams struct {
			Url string `json:"url"`
		}
		var url OpenBrowserParams
		err = event.DataAs(&url)
		if err != nil {
			svc.Logger().WithError(err).Warn("decode OpenBrowserParams")
			return nil, err
		}
		err = svc.OpenBrowser(url.Url)
		if err != nil {
			svc.Logger().WithError(err).Warn("browser open url")
			return nil, err
		}

	default:
		svc.Logger().WithFields(logrus.Fields{
			"topic": topic,
			"type":  event.Type(),
			"event": event,
		}).Warn("unknown event type")
	}

	return nil, nil
}
func (svc BrowserService) Run(ctx context.Context, params ...interface{}) error {
	log := svc.Logger()

	log.Debug("BrowserService started")
	defer log.Debug("BrowserService ended")

	var err error

	// register eventHandler for event reception
	err = svc.Transport().RegisterHandler(svc.eventHandler, "com.plugis.browser")
	if err != nil {
		svc.Logger().WithError(err).Error("failed to register event handler")
		return err
	}

	<-ctx.Done()

	return nil
}
