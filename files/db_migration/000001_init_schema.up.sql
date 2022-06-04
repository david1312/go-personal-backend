BEGIN;

create table customers
(
    id              int primary key AUTO_INCREMENT,
    uid             varchar(36) not null UNIQUE,
    name            varchar(100) not null,
    password        varchar(200) not null,
    email           varchar(200) not null,
    email_verified_token varchar(64),
    email_verified_at timestamp NULL DEFAULT NULL,
    email_verified_sent   tinyint NOT NULL default 0,
    email_change_code  varchar(6) NULL DEFAULT null,
    email_change_eligible      boolean not null default false,
    gender          enum('LAKI-LAKI', 'PEREMPUAN'),
    is_active       boolean not null default true,
    phone           varchar(20),
    phone_verified_at timestamp NULL DEFAULT NULL,
    avatar          varchar(100),
    birthdate       DATE NULL DEFAULT NULL,
    created_at      timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      timestamp NULL DEFAULT NULL
);

create table carts
(
    id              int primary key AUTO_INCREMENT,
    uid             varchar(36) not null UNIQUE,
    name            varchar(100) not null,
    password        varchar(200) not null,
    email           varchar(200) not null,
    email_verified_token varchar(64),
    email_verified_at timestamp NULL DEFAULT NULL,
    email_verified_sent   tinyint NOT NULL default 0,
    email_change_code  varchar(6) NULL DEFAULT null,
    email_change_eligible      boolean not null default false,
    gender          enum('LAKI-LAKI', 'PEREMPUAN'),
    is_active       boolean not null default true,
    phone           varchar(20),
    phone_verified_at timestamp NULL DEFAULT NULL,
    avatar          varchar(100),
    birthdate       DATE NULL DEFAULT NULL,
    created_at      timestamp DEFAULT CURRENT_TIMESTAMP,
    updated_at      timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at      timestamp NULL DEFAULT NULL
);

-- create type priority_events_enum as enum ('LOW','MEDIUM','HIGH','CRITICAL');
-- create table events
-- (
--     id                              serial primary key,
--     client_id                       integer not null,
--     name                            varchar(50) not null,
--     description                     varchar(100) not null,
--     priority                        priority_events_enum not null,
--     is_active                       boolean not null default false,

--     deleted_at                      timestamp with time zone,
--     created_at                      timestamp with time zone not null default current_timestamp,
--     updated_at                      timestamp with time zone,
--     created_by                      varchar(100),

--     foreign key (client_id) references clients (id)
-- );

-- create type type_templates_enum as enum ('WHATSAPP','EMAIL','SMS');
-- create table templates
-- (
--     id                      serial primary key,
--     event_id                integer not null,
--     name                    varchar(50) not null,
--     description             varchar(100),
--     content                 text,
--     status                  boolean not null default false,
--     external_template_id     varchar(50),

--     created_by  varchar(100) not null,

--     type        type_templates_enum not null,
--     priority    smallint not null,

--     deleted_at  timestamp with time zone,
--     created_at  timestamp with time zone not null default current_timestamp,
--     updated_at  timestamp with time zone,

--     foreign key (event_id) references events (id)
-- );

-- create table parameters
-- (
--     id              serial primary key,
--     event_id        integer not null,
--     name            varchar(50) not null,
--     is_required     boolean not null default false,

--     foreign key (event_id) references events (id)
-- );

-- create table incoming_requests
-- (
--     id                  bigserial primary key,
--     event_id            integer not null,
--     client_id           integer not null,
--     request_id          varchar(50) not null,
--     master_id           varchar(50),
--     status              boolean not null default false,
--     schedule            timestamp with time zone,
--     created_at          timestamp with time zone not null default current_timestamp,

--     foreign key (event_id) references events (id),
--     foreign key (client_id) references clients (id)
-- );

-- create table request_params
-- (
--     incoming_requests_id    int not null,
--     name   varchar(50)      not null,
--     value  varchar(100)     not null,
--     foreign key (incoming_requests_id) references incoming_requests (id)
-- );

-- create type type_activities_status_enum as enum ('SUCCESS','SENDING','FAILED');
-- create table activities
-- (
--     id                      bigserial primary key,
--     incoming_requests_id    integer not null,
--     template_id             integer not null,
--     status                  type_activities_status_enum,
--     status_msg              text,
--     interaction_id          varchar(50),
--     created_at              timestamp with time zone not null default current_timestamp,
--     log_msg                 text default '',
--     dispatch_status         text,

--     foreign key (incoming_requests_id) references incoming_requests (id),
--     foreign key (template_id) references templates (id)
-- );

COMMIT;
