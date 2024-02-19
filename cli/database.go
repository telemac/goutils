package cli

type Database struct {
	DbHost   string `help:"database host" default:"127.0.0.1" json:"db-host"`
	Database string `help:"database" json:"database"`
	DbUser   string `help:"database user" json:"db-user"`
	DbPass   string `help:"database password" json:"db-pass"`
	DbPort   int    `help:"database port" min:"0" max:"65535" json:"db-port"`
}
