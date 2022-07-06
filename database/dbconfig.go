package database

import (
	"fmt"
	"time"

	"github.com/jinzhu/configor"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBConfig struct {
	Host     string `default:"127.0.0.1" json:"host"`
	Database string `required:"true" json:"database"`
	User     string `default:"root" json:"user"`
	Password string `required:"true" env:"MysqlPassword"`
	Port     uint   `default:"3306" json:"port"`
}

// MySqlDSN returns the connection string for MySQL
func (dc *DBConfig) MySqlDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dc.User,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.Database,
	)
}

// MySQLConnect connects to a MySQL database
func MySQLConnect(config *DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(config.MySqlDSN()), &gorm.Config{})

	// avoid closing bad idle connection errors
	// TODO : check if needed with gorm 2
	d, err := db.DB()
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(100)
	d.SetMaxIdleConns(100)
	d.SetConnMaxLifetime(180 * time.Second)

	return db, err
}

// LoadMySQLConfig returns the mysql connection config
func LoadMySQLConfig(configFiles ...string) (*DBConfig, error) {

	type mySQL struct {
		Mysql DBConfig
	}
	var mySQLConfig mySQL

	err := configor.Load(&mySQLConfig, configFiles...)
	return &mySQLConfig.Mysql, err
}
