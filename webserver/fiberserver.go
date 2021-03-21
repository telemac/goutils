package webserver

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"sync"
)

type TemplateDataProviderFn func(c *fiber.Ctx) (fiber.Map, error)

type FiberServer struct {
	App                   *fiber.App
	viewsFolder           string
	staticFolder          string
	listenPort            int
	TemplateDataProviders map[string]TemplateDataProviderFn
	mutex                 sync.RWMutex
}

func NewFiberServer(viewsFolder string, staticFolder string, listenPort int) *FiberServer {
	engine := html.NewFileSystem(http.Dir("."+viewsFolder), ".html")
	engine.Reload(true)       // Optional. Default: false
	engine.Debug(true)        // Optional. Default: false
	engine.Layout("embed")    // Optional. Default: "embed"
	engine.Delims("{{", "}}") // Optional. Default: engine delimiters
	engine.AddFunc("greet", func(name string) string {
		return "Hello, " + name + "!"
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	server := &FiberServer{
		App:                   app,
		viewsFolder:           viewsFolder,
		staticFolder:          staticFolder,
		listenPort:            listenPort,
		TemplateDataProviders: make(map[string]TemplateDataProviderFn),
	}

	app.Static("/", "."+staticFolder)

	app.Get(viewsFolder+"*", func(c *fiber.Ctx) error {
		path := c.Path()
		path = strings.TrimPrefix(path, viewsFolder)
		if path == "" {
			path = "index"
		} else if strings.HasSuffix(path, "/") {
			path += "index"
		} else if strings.HasSuffix(path, ".html") {
			path = strings.TrimSuffix(path, ".html")
		}

		data := fiber.Map{
			"Title":  "Title from go",
			"Path":   path,
			"Colors": []string{"red", "green", "blue"},
		}

		// get datas for template via registred data providers
		dataProvider := server.GetTemplateDataProvider(path)
		if dataProvider != nil {
			dataMap, err := dataProvider(c)
			if err != nil {
				return fmt.Errorf("template %s data provider error : %w", path, err)
			}
			for k, v := range dataMap {
				data[k] = v
			}
		}

		log.WithFields(log.Fields{
			"path": path,
		}).Trace("get template")

		err := c.Render(path, data)
		if err != nil && strings.Contains(err.Error(), "does not exist") {
			data["Title"] = "Template not found"
			err = c.Render("404", data)
		}
		return err
	})
	return server
}

func (s *FiberServer) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		log.Warn("interrupted")
		err := s.App.Shutdown()
		if err != nil {
			log.WithError(err).Error("app.Shutdown")
		}
	}()
	return s.App.Listen(":3000")
}

func (s *FiberServer) AddTemplateDataProvider(path string, callback TemplateDataProviderFn) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	_, found := s.TemplateDataProviders[path]
	if found {
		return fmt.Errorf("template data profiver for path %s already set", path)
	}
	s.TemplateDataProviders[path] = callback
	return nil
}

func (s *FiberServer) GetTemplateDataProvider(path string) TemplateDataProviderFn {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.TemplateDataProviders[path]
}
