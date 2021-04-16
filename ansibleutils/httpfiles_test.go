package ansibleutils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHttpFiles_GetFile(t *testing.T) {
	assert := assert.New(t)
	httpFiles := NewHttpFiles("/tmp")
	localFile, err := httpFiles.GetFile("http://update.plugis.com/playbooks/site.yml")
	assert.NoError(err)
	assert.Equal("/tmp/site.yml", localFile)
}
