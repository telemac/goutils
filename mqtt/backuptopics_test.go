package mqtt

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBackupTopics(t *testing.T) {
	assert := assert.New(t)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*2)
	defer cancel()

	topics := []string{"colorbeam/loads",
		"colorbeam/building",
		"colorbeam/drivers",
		"colorbeam/calibrations",
		"colorbeam/persist",
	}

	topicsBackup, err := BackupTopics(ctx, "tcp://colorbeam:1883", topics)
	assert.NoError(err)
	_ = topicsBackup

}
