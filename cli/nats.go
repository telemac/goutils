package cli

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// NatsUrls holds the URLs of the NATS servers.
type NatsUrls struct {
	//Servers []string `help:"NATS server URLs." short:"s" default:"wss://user:pass@nats.domain.com:443"`
	Servers []string `help:"NATS server URLs." short:"s" default:"nats://localhost:4222"`
}

// NatsUser holds the user to connect to a NATS server.
type NatsUser struct {
	User string `help:"NATS server user." short:"u" default:""`
}

// NatsPass holds the password to connect to a NATS server.
type NatsPass struct {
	Pass string `help:"NATS server password." short:"p" default:""`
}

// NatsContext holds the context to connect to a NATS server, (empty string for default context	)
type NatsContext struct {
	Context string `help:"NATS context." short:"c" default:""`
}

// NatsConfig is a configuration structure for NATS.
type NatsConfig struct {
	NatsUrls
	NatsUser
	NatsPass
}

// NewNatsConfig creates a new NatsConfig.
func NewNatsConfig(natsUrls []string, natsUser, natsPass string) *NatsConfig {
	return &NatsConfig{
		NatsUrls: NatsUrls{
			Servers: natsUrls,
		},
		NatsUser: NatsUser{
			User: natsUser,
		},
		NatsPass: NatsPass{
			Pass: natsPass,
		},
	}
}

// NatsConnectString returns a connection string to be used to connect to a NATS server.
func (n *NatsConfig) NatsConnectString(serverIndex int) (string, error) {
	if serverIndex >= len(n.Servers) {
		return "", fmt.Errorf("invalid server index %d", serverIndex)
	}
	connectString, err := BuildConnectString(n.Servers[serverIndex], n.User, n.Pass)
	return connectString, err
}

// BuildConnectString builds a valid nats connection string given an url, user and passward
func BuildConnectString(urlStr, userStr, passwordStr string) (string, error) {
	// check if url contains a protocol prefix
	if !strings.Contains(urlStr, "://") {
		urlStr = "nats://" + urlStr
	}

	u, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	scheme := u.Scheme
	if scheme == "" {
		return "", errors.New("no scheme in url")
	}
	host := u.Host

	if host == "" && u.Path != "" {
		host = u.Path
	}

	user := u.User.Username()
	password, _ := u.User.Password()

	natsUrl := scheme + "://"

	if userStr != "" && passwordStr != "" {
		natsUrl += userStr + ":" + passwordStr + "@"
	} else if userStr != "" {
		natsUrl += userStr + "@"
	} else if passwordStr != "" {
		natsUrl += passwordStr + "@"
	} else if user != "" && password != "" {
		natsUrl += user + ":" + password + "@"
	} else if user != "" {
		natsUrl += user + "@"
	}

	natsUrl += host

	return natsUrl, nil
}
