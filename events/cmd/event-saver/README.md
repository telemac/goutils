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
