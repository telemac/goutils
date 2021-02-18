package heartbeat

import (
	"fmt"
	"github.com/telemac/goutils/net"
	"os"
	"time"
)

// ce type : com.plugis.heartbeat.Sent
type Sent struct {
	Hostname        string    `json:"hostname"`
	Mac             string    `json:"mac"`
	InternalIP      string    `json:"ip"`
	Started         time.Time `json:"started,omitempty"`
	Uptime          uint64    `json:"uptime,omitempty"`
	NatsServiceName string    `json:"nats-service,omitempty"`
}

// NewSent creates a new sent event
func NewSent(natsServiceName string) (*Sent, error) {
	// get external ip
	internalIP, err := net.GetOutboundIP()
	if err != nil {
		return nil, fmt.Errorf("get outbound ip :%w", err)
	}

	macAddress, err := net.GetMACAddress()
	if err != nil {
		return nil, fmt.Errorf("get mac address :%w", err)
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("get host name :%w", err)
	}

	return &Sent{
		Hostname:        hostname,
		Mac:             macAddress,
		InternalIP:      internalIP,
		Started:         time.Now(),
		NatsServiceName: natsServiceName,
	}, nil
}
