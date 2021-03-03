package cloudevents

import (
	"fmt"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

type Database struct {
	db *gorm.DB
}

type DatabaseConfig struct {
	DBHost string `ini:"dbhost,omitempty"`
	DBname string `ini:"dbname,omitempty"`
	DBuser string `ini:"dbuser,omitempty"`
	DBpass string `ini:"dbpass,omitempty"`
	DBPort int    `ini:"dbport,omitempty"`
}

func (d *Database) Open(dbConfig DatabaseConfig) error {
	// make logger dsn
	dsn := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable TimeZone=Europe/Paris",
		dbConfig.DBuser,
		dbConfig.DBpass,
		dbConfig.DBHost,
		dbConfig.DBPort,
		dbConfig.DBname,
	)

	var err error
	d.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	const createTableSql = `CREATE TABLE IF NOT EXISTS public.cloudevents (
                                    id UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
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

func (d *Database) InsertEvent(topic string, e *event.Event, payload []byte, err error) error {
	if err != nil { // malformed event
		const insertSql = `insert into public.cloudevents (id, time, type, topic) values (?,?,?,'{}')`
		id := uuid.NewString()
		t := time.Now()
		eventType := "malformed"
		return d.db.Exec(insertSql, id, t, eventType, topic).Error
	}
	const insertSql = `insert into public.cloudevents (id, time, type, topic, data, datacontenttype, source, specversion) values (?,?,?,?,?,?,?,?)`
	return d.db.Exec(insertSql, e.ID(), e.Time(), e.Type(), topic, e.Data(), e.DataContentType(), e.Source(), e.SpecVersion()).Error
}
