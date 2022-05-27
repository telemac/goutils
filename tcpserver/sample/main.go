package main

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/task"
	"github.com/telemac/goutils/tcpserver"
)

type MyConnectionHandler struct {
}

func (c *MyConnectionHandler) OnConnect(ctx context.Context, connection tcpserver.ConnectionParams) bool {
	connection.ServerConfig.Logger.WithField("remote", connection.Conn.RemoteAddr().String()).Info("OnConnect")
	return true
}

func (c *MyConnectionHandler) OnDisconnect(ctx context.Context, connection tcpserver.ConnectionParams) {
	connection.ServerConfig.Logger.WithField("remote", connection.Conn.RemoteAddr().String()).Info("OnDisconnect")
}

func (c *MyConnectionHandler) HandleConnection(ctx context.Context, connection tcpserver.ConnectionParams) error {
	connection.ServerConfig.Logger.WithField("remote", connection.Conn.RemoteAddr().String()).Info("HandleConnection")
	buffer := make([]byte, 1024)
	connection.Conn.Write([]byte("Ok\n"))

	go func() {
		select {
		case <-ctx.Done():
			connection.Conn.Close()
		}
	}()

	for {
		n, err := connection.Conn.Read(buffer)
		connection.ServerConfig.Logger.WithFields(logrus.Fields{
			"remote":   connection.Conn.RemoteAddr().String(),
			"read":     string(buffer[:n]),
			"UserData": connection.ServerConfig.UserData,
			"error":    err,
		}).Info("Read")
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {

	ctx, cancel, logger := task.NewCancellableContextWithLog(time.Second*5, "trace", logrus.Fields{
		"app": "tcp server demo",
	})
	defer cancel()

	server := tcpserver.NewServer(tcpserver.ServerConfig{
		ListenPort: 8080,
		UserData:   "My user data", // can be anything, must be thread safe
		Logger:     logger,
		Connection: &MyConnectionHandler{},
	})

	services := task.NewRunnerRepository()
	services.Start(ctx, server)
	services.WaitUntilAllDone()

	//err := server.Run(ctx)
	//if err != nil {
	//	logger.WithError(err).Error("run tcp server")
	//}

}
