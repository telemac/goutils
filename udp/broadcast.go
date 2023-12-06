package udp

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"

	reuse "github.com/libp2p/go-reuseport"
)

// UDPBroadcaster allow to boradcast and receive broadcasted udp datagrams
type UDPBroadcaster struct {
	listener net.PacketConn
	port     int
	buffer   []byte
}

type UDPDatagram struct {
	Payload []byte
	Addr    net.Addr
}

func NewUDPBroadcaster(port int) (*UDPBroadcaster, error) {
	udpBroadcaster := &UDPBroadcaster{
		port:   port,
		buffer: make([]byte, 1500),
	}
	var err error
	udpBroadcaster.listener, err = reuse.ListenPacket("udp", fmt.Sprintf(":%d", udpBroadcaster.port))
	return udpBroadcaster, err
}

func (ub *UDPBroadcaster) Close() {
	if ub.listener != nil {
		ub.listener.Close()
	}
}

func (ub *UDPBroadcaster) Broadcast(databram []byte) error {
	if ub.listener == nil {
		return errors.New("listener not initialized")
	}
	addr, err := reuse.ResolveAddr("udp", "255.255.255.255:"+strconv.Itoa(ub.port))
	if err != nil {
		return err
	}
	n, err := ub.listener.WriteTo(databram, addr)
	if err != nil {
		return err
	}
	if n != len(databram) {
		return errors.New("partial send")
	}
	return nil
}

// Unicast sends an udp datagram to a specific address
func (ub *UDPBroadcaster) Unicast(destAddr string, databram []byte) error {
	if ub.listener == nil {
		return errors.New("listener not initialized")
	}
	addr, err := reuse.ResolveAddr("udp", destAddr+":"+strconv.Itoa(ub.port))

	if err != nil {
		return err
	}
	n, err := ub.listener.WriteTo(databram, addr)
	if err != nil {
		return err
	}
	if n != len(databram) {
		return errors.New("partial send")
	}
	return nil
}

func (ub *UDPBroadcaster) Read(timeout time.Duration) (*UDPDatagram, error) {
	err := ub.listener.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return nil, err
	}
	n, addr, err := ub.listener.ReadFrom(ub.buffer)
	if n > 0 {
		datagram := new(UDPDatagram)
		datagram.Addr = addr
		datagram.Payload = make([]byte, n)
		copy(datagram.Payload, ub.buffer[:n])
		return datagram, err
	}
	return nil, err
}

// Broadcast broadcasts an udp datagram
func Broadcast(addrAndPort string, buffer []byte) error {
	c, err := reuse.Dial("udp", "", addrAndPort)
	if err != nil {
		return err
	}
	n, err := c.Write(buffer)
	if n != len(buffer) {
		return errors.New("partial send")
	}
	return err
}
