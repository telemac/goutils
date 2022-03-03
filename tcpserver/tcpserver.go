package tcpserver

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net"
	"sync"
	"sync/atomic"

	"github.com/tevino/abool"
)

var (
	ErrCancelled = errors.New("cancelled")
)

// TCPServer holds the data for our TCP Server
type TCPServer struct {
	wg         sync.WaitGroup // waitGroup for running goroutines
	accepted   int64          // Nb of accepted connexions
	goroutines int64          // Nb active goroutines handling connections
	tlsConfig  *tls.Config    // tlsConfig parameters, TLS if set
}

// ConnectionHandler is the signature of the handler function
// Returning an error in the handler exits the ListerAndServe function with that error
type ConnectionHandler func(ctx context.Context, conn net.Conn, userData interface{})

// Wait waits until all the goroutines are finished
func (tcpServer *TCPServer) Wait() {
	tcpServer.wg.Wait()
}

// GetNbAccepted returns the number of connections accepted
func (tcpServer *TCPServer) GetNbAccepted() int64 {
	return atomic.LoadInt64(&tcpServer.accepted)
}

// GetNbGoroutines returns the number of goroutines currently handling connections
func (tcpServer *TCPServer) GetNbGoroutines() int64 {
	return atomic.LoadInt64(&tcpServer.goroutines)
}

// ListerAndServe accepts incoming tcp connexions and calls the handler in a new goroutine
func (tcpServer *TCPServer) ListerAndServe(ctx context.Context, address string, handler ConnectionHandler, tlsConfig *tls.Config, userData interface{}) error {
	tcpServer.tlsConfig = tlsConfig

	// Listen for incoming tcp connections
	// Use a cancellable listen
	var listenConfig net.ListenConfig
	tcpNetListener, err := listenConfig.Listen(ctx, "tcp", address)
	if err != nil {
		return err
	}
	defer tcpNetListener.Close()

	var netListener = tcpNetListener

	// Is TLS
	if tcpServer.tlsConfig != nil {
		tlsNetListener := tls.NewListener(tcpNetListener, tcpServer.tlsConfig)
		netListener = tlsNetListener
		defer tlsNetListener.Close()
	}

	log.Printf("Listening on %s", address)

	interrupted := abool.New() // true if context is cancelled

	// Handle context cancellation
	go func(interrupted *abool.AtomicBool) {
		<-ctx.Done()
		interrupted.Set()
		netListener.Close()
	}(interrupted)

	// Loop on incoming connections
	for {
		conn, err := netListener.Accept()
		if err != nil {
			if interrupted.IsSet() {
				return ErrCancelled
			}
			return err
		}
		atomic.AddInt64(&tcpServer.accepted, 1)
		tcpServer.wg.Add(1) // used for TCPServer.Wait()
		go func(tcpServer *TCPServer, conn net.Conn) {
			atomic.AddInt64(&tcpServer.goroutines, 1) // count active goroutines
			handler(ctx, conn, userData)
			conn.Close() // May be closed twice, better than not.
			tcpServer.wg.Done()
			atomic.AddInt64(&tcpServer.goroutines, -1)
		}(tcpServer, conn)
	}
}
