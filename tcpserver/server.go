package tcpserver

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"strconv"

	"github.com/sirupsen/logrus"
	"github.com/telemac/goutils/task"
)

// ServerConfig holds the server configuration parameters for a Server
type ServerConfig struct {
	ListenAddress string // "0.0.0.0"
	ListenPort    int    // 8080
	TlsConfig     *tls.Config
	Logger        *logrus.Entry
	UserData      interface{} // any datas, passed to each connection, must be thread safe
	Connection    ConnectionIntf
}

// Server represents a tcp server
type Server struct {
	task.Runnable
	config ServerConfig
}

func NewServer(config ServerConfig) *Server {
	return &Server{config: config}
}

func (s Server) Run(ctx context.Context, params ...interface{}) error {
	listenStr := s.config.ListenAddress + ":" + strconv.Itoa(s.config.ListenPort)

	// create standard logger if not defined
	if s.config.Logger == nil {
		s.config.Logger = logrus.New().WithField("logrus", "default logger")
		s.config.Logger.Warn("no logger defined for tcp server")
	}

	// add listen field to the logger
	s.config.Logger = s.config.Logger.WithFields(logrus.Fields{
		"listen_on": listenStr,
	})

	var tcpServer TCPServer
	err := tcpServer.ListerAndServe(ctx, listenStr, s.connectionHandler, s.config.TlsConfig, s.config.UserData)
	if err != nil {
		if errors.Is(err, ErrCancelled) {
			s.config.Logger.WithError(err).Warn("listen to incoming connections cancelled")
		} else {
			s.config.Logger.WithError(err).Error("listen to incoming connections")
		}
	}
	tcpServer.Wait()

	return nil

}

func (s Server) connectionHandler(ctx context.Context, conn net.Conn, userData interface{}) {
	// one goroutine per connection
	connectionParams := ConnectionParams{
		Conn:         conn,
		ServerConfig: &s.config,
	}
	acceptConnection := connectionParams.ServerConfig.Connection.OnConnect(ctx, connectionParams)
	if acceptConnection {
		connectionParams.ServerConfig.Connection.HandleConnection(ctx, connectionParams)
	}
	connectionParams.ServerConfig.Connection.OnDisconnect(ctx, connectionParams)

}
