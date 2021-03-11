package heartbeat

import (
	"context"
	"errors"
	"github.com/foolin/goview"
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/telemac/goutils/natsservice"
	"net/http"
	"time"
)

// HeartbeatWebInterface exposes com.plugis.heartbeat.Sent events
// that are saved in the database
type HeartbeatWebInterface struct {
	natsservice.NatsService
	db Database
}

func (svc *HeartbeatWebInterface) getIndex(ctx *gin.Context) {
	//render with master
	records, err := svc.db.getHeartbeat(true)
	ctx.HTML(http.StatusOK, "index", gin.H{
		"title":      "Index title!",
		"name":       "Alexandre",
		"heartbeats": records,
		"error":      err,
		"add": func(a int, b int) int {
			return a + b
		},
	})
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

	router := gin.Default()

	//new template engine
	config := goview.DefaultConfig
	config.DisableCache = true
	router.HTMLRender = ginview.New(config)

	router.GET("/", svc.getIndex)

	router.GET("/page", func(ctx *gin.Context) {
		//render only file, must full name with extension
		ctx.HTML(http.StatusOK, "page.html", gin.H{"title": "Page file title!!"})
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error("listen")
			// TODO : exit Run function on error
		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.WithError(err).Error("Server Shutdown")
		return err
	}

	return nil
}
