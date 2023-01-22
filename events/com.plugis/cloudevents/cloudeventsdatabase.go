package cloudevents

import (
	"fmt"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"github.com/telemac/goutils/natsservice"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CloudEventsDatabase struct {
	db *gorm.DB
}

func (d *CloudEventsDatabase) Open(dbConfig natsservice.PostgresConfig) error {
	// make logger dsn
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable TimeZone=Europe/Paris",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database,
	)

	var err error
	d.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	const createExtension = `CREATE EXTENSION if not exists "uuid-ossp";`
	err = d.db.Exec(createExtension).Error
	if err != nil {
		return err
	}

	const createTableSql = `CREATE TABLE IF NOT EXISTS cloudevents (
                                    id UUID DEFAULT uuid_generate_v4()::UUID PRIMARY KEY,
                                    time timestamp not null ,
                                    type VARCHAR NOT NULL,
                                    topic varchar not null,
                                    data JSONB NULL,
                                    datacontenttype varchar,
                                    source varchar,
                                    specversion varchar
);`
	return d.db.Exec(createTableSql).Error
}

// Cleanheartbeats keeps the last occurence for each heartbeat
func (d *CloudEventsDatabase) Cleanheartbeats() error {
	const cleanHeartbeatsSql = `delete from cloudevents where id in (
    select id from (
                       select topic,max(time) as time from cloudevents
                       where topic like 'com.plugis.heartbeat.Sent.%'
                       group by topic
                   ) as lastHeartbeats,cloudevents
    where cloudevents.topic=lastHeartbeats.topic and cloudevents.time < lastHeartbeats.time
    )`
	return d.db.Exec(cleanHeartbeatsSql).Error
}

func (d *CloudEventsDatabase) InsertEvent(topic string, e *event.Event, payload []byte, err error) error {
	if err != nil { // malformed event
		const insertSql = `insert into cloudevents (id, time, type, topic) values (?,?,?,'{}')`
		id := uuid.NewString()
		t := time.Now()
		eventType := "malformed"
		return d.db.Exec(insertSql, id, t, eventType, topic, payload).Error
	}
	if e.Type() == "com.plugis.heartbeat.Sent" {
		// delete last heartbeat events here
		const deleteHeartbeatSql = `delete from cloudevents where topic=?`
		err = d.db.Exec(deleteHeartbeatSql, topic).Error
		if err != nil {
			return err
		}
	}
	const insertSql = `insert into cloudevents (id, time, type, topic, data, datacontenttype, source, specversion) values (?,?,?,?,?,?,?,?)`
	return d.db.Exec(insertSql, e.ID(), e.Time(), e.Type(), topic, e.Data(), e.DataContentType(), e.Source(), e.SpecVersion()).Error
}
