package database

import (
	"errors"
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUnknownDBType = errors.New("unknown database type")
	ErrDBNameEmpty   = errors.New("database name is empty")
	ErrDBHostEmpty   = errors.New("database host is empty")
)

type DBNameType string

// GenericDBConfigMap holds the database configs by name
type GenericDBConfigMap map[DBNameType]GenericDBConfig

// GenericDBConfig holds database config for MySQL/Postgre/...
type GenericDBConfig struct {
	Name     DBNameType `required:"true" json:"name"` // database connection name, ex : orders-mysql
	Type     string     `required:"true" json:"type"` // mysql/postgre
	Host     string     `required:"true" json:"host"`
	Database string     `required:"true" json:"database"`
	User     string     `required:"true" json:"user"`
	Password string     `required:"true"`
	Port     uint       `required:"true" json:"port"`
}

func (dbConfig *GenericDBConfig) Validate() error {
	if dbConfig.Type != "mysql" && dbConfig.Type != "postgre" {
		return ErrUnknownDBType
	}
	if dbConfig.Name == "" {
		return ErrDBNameEmpty
	}
	if dbConfig.Host == "" {
		return ErrDBHostEmpty
	}
	// TODO : add more checks
	return nil
}

// DSN returns the connection string corresponding to the database type
func (dbConfig *GenericDBConfig) DSN() (string, error) {
	switch dbConfig.Type {
	case "mysql":
		return fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Database,
		), nil
	case "postgre":
		return fmt.Sprintf(
			"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable TimeZone=Europe/Paris",
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Host,
			dbConfig.Port,
			dbConfig.Database,
		), nil
	default:
		return "", fmt.Errorf("unknown database type %s : %w", dbConfig.Type, ErrUnknownDBType)
	}
}

// Connect connects to the database
func (dbConfig *GenericDBConfig) Connect() (*gorm.DB, error) {
	var db *gorm.DB
	dsn, err := dbConfig.DSN()
	if err != nil {
		return nil, err
	}
	switch dbConfig.Type {
	case "mysql":
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	case "postgre":
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	default:
		return nil, fmt.Errorf("unknown database type %s : %w", dbConfig.Type, ErrUnknownDBType)
	}
	// check gorm.Open error for the corresponding database type
	if err != nil {
		return nil, fmt.Errorf("can not open database dsn=%s : %w", dsn, err)
	}

	// avoid closing bad idle connection errors
	// TODO : check if needed with gorm 2
	d, err := db.DB()
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(100)
	d.SetMaxIdleConns(100)
	d.SetConnMaxLifetime(10 * time.Second)

	return db, nil
}

/*
LoadGenericDBConfig reads a config yaml file
Sapple config file :
databases:
  - type: mysql
    name: my-mysql-db
    host: myMysqlHost
    database: myMysqlDatabase
    user:     myMysqlUser
    password: myMysqlPassword
    port:     3306
  - type: postgre
    name: my-postgres-db
    host: myPostgresHost
    database: myPostgresDatabase
    user:     myPostgresUser
    password: myPostgresPassword
    port:     5432
*/
func LoadGenericDBConfig(fileName string, genericDBConfigMap *GenericDBConfigMap) error {
	if *genericDBConfigMap == nil {
		*genericDBConfigMap = make(GenericDBConfigMap)
	}
	// load config from file
	var k = koanf.New(".")
	err := k.Load(file.Provider(fileName), yaml.Parser())
	if err != nil {
		return err
	}
	// read database configs as array
	var dbConfigArray []GenericDBConfig
	err = k.Unmarshal("databases", &dbConfigArray)
	if err != nil {
		return err
	}
	// convert to map
	for _, dbConfig := range dbConfigArray {
		err = dbConfig.Validate()
		if err != nil {
			return err
		}
		(*genericDBConfigMap)[dbConfig.Name] = dbConfig
	}
	return nil
}
