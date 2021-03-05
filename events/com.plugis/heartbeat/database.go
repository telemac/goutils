package heartbeat

import (
	"fmt"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	DBHost string `ini:"dbhost,omitempty"`
	DBname string `ini:"dbname,omitempty"`
	DBuser string `ini:"dbuser,omitempty"`
	DBpass string `ini:"dbpass,omitempty"`
	DBPort int    `ini:"dbport,omitempty"`
}

type Database struct {
	db *gorm.DB
}

func (d *Database) Open(dbConfig DatabaseConfig) error {
	// make logger dsn
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.DBuser,
		dbConfig.DBpass,
		dbConfig.DBHost,
		dbConfig.DBPort,
		dbConfig.DBname,
	)

	var err error
	d.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	const createTableSql = `create table if not exists heartbeats
(
    id              int auto_increment
        primary key,
    mac             varchar(17) not null,
    hostname        varchar(64) null,
    ip              varchar(64) null,
    first_heartbeat datetime    null,
    last_heartbeat  datetime    not null,
    comment         tinytext    null,
    constraint heartbeats_mac_uindex
        unique (mac)
);`
	return d.db.Exec(createTableSql).Error
}

func (d *Database) upsertHeartbeat(sent Sent) error {

	const upsertSql = `insert into heartbeats
    (mac, last_heartbeat, first_heartbeat, hostname,ip)
VALUES
    (?,CURRENT_TIMESTAMP(),CURRENT_TIMESTAMP(), ?, ?)
on duplicate key update
    last_heartbeat=CURRENT_TIMESTAMP(),ip=?;`
	err := d.db.Exec(upsertSql, sent.Mac, sent.Hostname, sent.InternalIP, sent.InternalIP).Error

	return err
}
