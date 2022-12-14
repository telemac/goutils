# cloud event saver

## create postgres database
```postgresql
create table cloudevents
(
    id              uuid default uuid_generate_v4() not null
        constraint cloudevents_pkey
            primary key,
    time            timestamp                       not null,
    type            varchar                         not null,
    topic           varchar                         not null,
    data            jsonb,
    datacontenttype varchar,
    source          varchar,
    specversion     varchar
);

alter table cloudevents
    owner to postgres;

create index time_idx
    on cloudevents (time);

grant delete, insert, references, select, trigger, truncate, update on cloudevents to plugis;
```

## cloud events request samples
```postgresql
-- get 100 newest cloud events
select * from cloudevents order by time desc limit 100;

-- cloud events per type/topic
select type,topic,count(*) from cloudevents
group by type,topic order by count desc;

-- SMS reçus extraction directe des could events
select TO_CHAR((time - cast(data->>'timestamp' as timestamptz)),'HH24:MI:SS') as difference,time,data->>'timestamp' as timestamp,data->>'from' as from,data->>'message' as message from cloudevents
where type='com.megalarm.sms.received'
order by data->>'timestamp' desc

-- delete all heartbeats
delete from cloudevents
where type='com.plugis.heartbeat.Sent'


```
