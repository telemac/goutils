package cloudevents

import (
	"fmt"
	"github.com/cloudevents/sdk-go/v2/types"
	"github.com/telemac/goutils/natsservice"
	"github.com/telemac/goutils/variable"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

type VariablesDatabase struct {
	db *gorm.DB
}

func (d *VariablesDatabase) Open(dbConfig natsservice.MysqlConfig) error {
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

	const createTableSql = `create table if not exists variables
(
    name        varchar(255)                        not null primary key,
    type        varchar(16)                         not null,
    value       longtext                            null,
    unit        varchar(64)                         null,
    comment     varchar(256)                        null,
    event_uuid  varchar(36)                         null,
    created_at  timestamp default CURRENT_TIMESTAMP not null,
    modified_at timestamp default CURRENT_TIMESTAMP not null,
    constraint variables_name_uindex
        unique (name)
);`
	return d.db.Exec(createTableSql).Error
}

func (d *VariablesDatabase) upsertVariables(event_uuid string, variables variable.Variables) error {
	const upsertSql = `insert into variables
    (name,type,value,unit,comment,event_uuid,created_at,modified_at)
VALUES
    (?,?,?,?,?,?,?,?)
on duplicate key update
    type=?, value=?, unit=?, comment=?, event_uuid=?, modified_at=?`
	for _, v := range variables {
		if types.IsZero(v.Timestamp) {
			v.Timestamp = time.Now()
		}
		err := d.db.Exec(upsertSql, v.Name, v.Type, v.Value, v.Unit, v.Comment, event_uuid, v.Timestamp, v.Timestamp, v.Type, v.Value, v.Unit, v.Comment, event_uuid, v.Timestamp).Error
		if err != nil {
			return fmt.Errorf("upsert variable %s : %w", v.Name, err)
		}
	}
	return nil
}
