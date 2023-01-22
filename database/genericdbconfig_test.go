package database

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadGenericDBConfig(t *testing.T) {
	assert := assert.New(t)
	var genericDBConfigMap GenericDBConfigMap
	err := LoadGenericDBConfig("./testdata/multiple-databases.yaml", &genericDBConfigMap)
	assert.NoError(err)
	assert.Equal(GenericDBConfigMap{"my-mysql-db": GenericDBConfig{Type: "mysql", Name: "my-mysql-db", Host: "myMysqlHost", Database: "myMysqlDatabase", User: "myMysqlUser", Password: "myMysqlPassword", Port: 0xcea}, "my-postgres-db": GenericDBConfig{Type: "postgres", Name: "my-postgres-db", Host: "myPostgresHost", Database: "myPostgresDatabase", User: "myPostgresUser", Password: "myPostgresPassword", Port: 0x1538}}, genericDBConfigMap)
}
