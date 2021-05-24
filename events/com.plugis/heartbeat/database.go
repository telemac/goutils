package heartbeat

import (
	"fmt"
	"github.com/telemac/goutils/natsservice"
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

func (d *Database) Open(dbConfig natsservice.MysqlConfig) error {
	// make logger dsn
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
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

func (d *Database) getHeartbeats(valid bool) ([]map[string]interface{}, error) {
	//	sqlStr := ``
	//	if valid {
	//		sqlStr = `select mac,TIMEDIFF(now(),last_heartbeat) as elapsed,last_heartbeat
	//from plugis.heartbeats
	//where TIMEDIFF(now(),last_heartbeat)<=60
	//order by last_heartbeat desc;`
	//	} else {
	//		sqlStr = `sselect mac,TIMEDIFF(now(),last_heartbeat) as elapsed,last_heartbeat
	//from plugis.heartbeats
	//where TIMEDIFF(now(),last_heartbeat)>60
	//order by last_heartbeat desc;`
	//	}
	//	tx := d.db.Exec(sqlStr)
	//	if tx.Error != nil {
	//		return nil, tx.Error
	//	}
	var dest []map[string]interface{}
	err := d.db.Table("heartbeats").Order("first_heartbeat desc").Select("*,TIMEDIFF(now(),last_heartbeat) as elapsed").Find(&dest).Error
	if err != nil {
		return nil, err
	}

	return dest, nil
}

/*
# valid heartbeats
select mac,TIMEDIFF(now(),last_heartbeat) as elapsed,last_heartbeat
from plugis.heartbeats
where TIMEDIFF(now(),last_heartbeat)<=60
order by last_heartbeat desc;

# invalid heartbeats
select mac,TIMEDIFF(now(),last_heartbeat) as elapsed,last_heartbeat
from plugis.heartbeats
where TIMEDIFF(now(),last_heartbeat)>60
order by last_heartbeat desc;
*/
