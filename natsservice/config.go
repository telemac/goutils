package natsservice

import (
	"github.com/jinzhu/configor"
)

type PostgresConfig struct {
	Host     string `default:"127.0.0.1"`
	Database string `default:"plugis"`
	User     string `default:"plugis"`
	Password string `required:"true" env:"PostgresPassword" default:"plugis"`
	Port     uint   `default:"5432"`
}

type MysqlConfig struct {
	Host     string `default:"127.0.0.1"`
	Database string `default:"plugis"`
	User     string `default:"plugis"`
	Password string `required:"true" env:"MysqlPassword" default:"plugis"`
	Port     uint   `default:"3306"`
}

type Config struct {
	//APPName string `default:"event-saver"`

	Servers []struct {
		Url string `default:"nats://nats1.plugis.com:443"`
		Sni bool   `default:"false"`
	}

	Postgres PostgresConfig

	Mysql MysqlConfig
}

func LoadConfig(files ...string) (Config, error) {
	var config Config
	err := configor.Load(&config, files...)
	return config, err
}
