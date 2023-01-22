package remote_access

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestRemoteAccess_Run(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	remoteAccess := NewRemoteAccess(RemoteAccessconfig{
		BaseUpdateUrl: "https://update.plugis.com/",
	})
	err := remoteAccess.Run(ctx)

	assert.NoError(err)
}
