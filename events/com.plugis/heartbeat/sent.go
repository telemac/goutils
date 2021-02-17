package heartbeat

import (
	"fmt"
	"github.com/telemac/goutils/net"
	"time"
)

// ce type : com.plugis.heartbeat.Sent
type Sent struct {
	Mac        string    `json:"mac"`
	InternalIP string    `json:"ip"`
	Started    time.Time `json:"started"`
}

// NewSent creates a new sent event
func NewSent() (*Sent, error) {
	// get external ip
	externelIP, err := net.GetOutboundIP()
	if err != nil {
		fmt.Errorf("get outbound ip :%w", err)
	}

	macAddress, err := net.GetMACAddress()
	if err != nil {
		fmt.Errorf("get mac address :%w", err)
	}

	return &Sent{
		Mac:        macAddress,
		InternalIP: externelIP,
		Started:    time.Now(),
	}, nil
}
