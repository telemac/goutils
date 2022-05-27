package tcpserver

import (
	"net"

	"golang.org/x/net/context"
)

type ConnectionParams struct {
	Conn         net.Conn
	ServerConfig *ServerConfig
}

type ConnectionIntf interface {
	OnConnect(ctx context.Context, connection ConnectionParams) bool
	OnDisconnect(ctx context.Context, connection ConnectionParams)
	HandleConnection(ctx context.Context, connection ConnectionParams) error
}
