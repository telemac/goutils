package natsservice

import (
	"flag"
	"fmt"
	"github.com/jinzhu/configor"
	"os"
	"path"
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
		Sni bool   `default:"true"`
	}

	Postgres PostgresConfig

	Mysql MysqlConfig

	CommandLineParams CommonCommandLineParams
}

// CommandLineParams holds the common, non generic command line parameters
type CommonCommandLineParams struct {
	Log          string
	ConfigFolder string
}

// ParseCommonCommandLineParams sets the command line parameters, or their default values
func ParseCommonCommandLineParams(params *CommonCommandLineParams) {
	flag.StringVar(&params.Log, "log", "info", "log level")
	flag.StringVar(&params.ConfigFolder, "config", "./", "folder containing configuration files")
	flag.Parse()
}

func LoadConfig(files ...string) (Config, error) {
	var config Config
	ParseCommonCommandLineParams(&config.CommandLineParams)

	// append config folder to given files
	var configFiles []string
	for _, file := range files {
		// check if file exists.
		fullPath := path.Join(config.CommandLineParams.ConfigFolder, file)
		_, err := os.Stat(fullPath)
		if err != nil {
			return config, fmt.Errorf("load config file %s : %w", fullPath, err)
		}
		configFiles = append(configFiles, fullPath)
	}

	err := configor.Load(&config, configFiles...)
	return config, err
}
