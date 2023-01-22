# Save cloud events to PostgreSQL

```postgresql
-- get postgres version
SELECT version();

-- add uuid extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- cloudevents stores a log for cloud events
CREATE TABLE IF NOT EXISTS public.cloudevents (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    time timestamp not null ,
    type VARCHAR NOT NULL,
    topic varchar not null,
    data JSONB NULL,
    datacontenttype varchar,
    source varchar,
    specversion varchar
);

-- create indexes
create index time_idx  on public.cloudevents(time);
create index type_idx  on public.cloudevents(type);
create index topic_idx  on public.cloudevents(topic);

-- create plugis user and give access rights
create user plugis with encrypted password 'plugis';
-- change password if needed
ALTER USER plugis WITH encrypted password 'plugis';

grant all privileges on database plugis to plugis;
grant all privileges on table plugis.public.cloudevents to plugis;

# count saved events
select count(*) from cloudevents;

# last 100 events
select * from cloudevents order by time desc limit 100;

# all topics
select distinct topic from cloudevents

select distinct data->'mac' as hostname,data->'started' as started from cloudevents where type='com.plugis.heartbeat.Sent'

-- show windspeed
select time,jsonb_array_element(data,1)->'value' as windspeed from cloudevents where topic='com.drone-box.box.1' and type='com.drone-box.weather.tempest.RapidWind'

-- select the first variable of the array, sample request
select time,name,value from (
                                                 select time,topic,jsonb_array_element(data,0)->'name'->>0 as name,jsonb_array_element(data,0)->'value'->>0 as value from public.cloudevents
                                                 where type='com.plugis.variable.set' and time > now() - interval '30 minute'
                                             ) as VARIABLES
where name='tempest.ST-00030095.solar-radiation' and value::integer < 600
order by time desc

-- potentiel solar production sample
select time,name,value::numeric*14.5,pg_typeof(value) from (
                                                               select time,topic,jsonb_array_element(data,0)->'name'->>0 as name,jsonb_array_element(data,0)->'value'->>0 as value from public.cloudevents
                                                               where type='com.plugis.variable.set' and time > now() - interval '300 minute'
                                                           ) as VARIABLES
where name='tempest.ST-00030095.solar-radiation'
order by time desc
limit 10

-- get last heartbeats
select data->'mac' as mac,max(time) as last,min(time) as first from cloudevents where type='com.plugis.heartbeat.Sent' group by mac
order by last desc

-- extract heartbeats
SELECT id, type,
       data->>'ip' AS ip,
       data->>'os' AS os,
       data->>'mac' AS mac,
       data->>'arch' AS arch,
       data->>'uptime' AS uptime,
       data->>'started' AS started,
       data->>'hostname' AS hostname,
       data->>'nats-service' AS nats_service
FROM cloudevents
where type='com.plugis.heartbeat.Sent'
order by time desc;

```


