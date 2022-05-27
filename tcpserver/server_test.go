package tcpserver

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type MyConnectionHandler struct {
}

func (c *MyConnectionHandler) OnConnect(ctx context.Context, connection ConnectionParams) bool {
	connection.ServerConfig.Logger.WithField("remote", connection.Conn.RemoteAddr().String()).Info("OnConnect")
	return true
}

func (c *MyConnectionHandler) OnDisconnect(ctx context.Context, connection ConnectionParams) {
	connection.ServerConfig.Logger.WithField("remote", connection.Conn.RemoteAddr().String()).Info("OnDisconnect")
}

func (c *MyConnectionHandler) HandleConnection(ctx context.Context, connection ConnectionParams) error {
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

func TestNewServer(t *testing.T) {
	assert := assert.New(t)

	var myConnectionHandler MyConnectionHandler

	server := NewServer(ServerConfig{
		ListenPort: 8080,
		UserData:   "My user data", // can be anything, must be thread safe
		Connection: &myConnectionHandler,
	})
	assert.NotNil(server)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := server.Run(ctx)
	assert.NoError(err)

}
