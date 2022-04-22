package database

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadMySQLConfig(t *testing.T) {
	assert := assert.New(t)

	configFile := "testdata/mysql.yaml"

	_, err := os.Stat(configFile)
	assert.NoError(err)

	got, err := LoadMySQLConfig(configFile)
	assert.NoError(err)
	assert.Equal(&DBConfig{Host: "myHost", Database: "myDatabase", User: "myUser", Password: "myPassword", Port: 3306}, got)
}
