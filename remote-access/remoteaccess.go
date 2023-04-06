package remote_access

import "github.com/telemac/goutils/natsservice"

// RemoteAccessconfig holds the configuration parameters for remote-access service
type Config struct {
	BaseUpdateUrl  string                        // BaseUpdateUrl where to download updates ex: https://update.plugis.com/ (must end with /)
	NatsServers    []string                      // list of nats servers, tries to connect from first to last
	NatsServices   []natsservice.NatsServiceIntf // nats services to launch
	HeartbeatMetas map[string]interface{}        // metas to send with heartbeat
}

// RemoteAccess holds data for a remote-access process
type RemoteAccess struct {
	commandLineParams CommandLineParams
	config            Config
}

func NewRemoteAccess(config Config) *RemoteAccess {
	return &RemoteAccess{config: config}
}
