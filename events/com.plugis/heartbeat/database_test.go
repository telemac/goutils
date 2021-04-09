package heartbeat

import (
	"github.com/stretchr/testify/assert"
	"github.com/telemac/goutils/natsservice"
	"testing"
)

func TestDatabase_Open(t *testing.T) {
	assert := assert.New(t)
	var db Database
	dbConfig := natsservice.MysqlConfig{
		Host:     "127.0.0.1",
		Database: "plugis",
		User:     "root",
		Password: "telemac",
		Port:     3306,
	}
	err := db.Open(dbConfig)
	assert.NoError(err)
	err = db.upsertHeartbeat(Sent{
		Mac: "11:22:33:44:55:66",
	})
	assert.NoError(err)
}
