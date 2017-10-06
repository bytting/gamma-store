drop database gs;
create database gs;

\c gs;

create sequence users_id;

create table users (
    id int not null default nextval('users_id'),
    username varchar not null,
    password varchar not null
);

alter sequence users_id owned by users.id;
grant all privileges on table users to numsys;
grant usage, select on sequence users_id to numsys;

create sequence spectrum_id;

create table spectrum (
    id int not null default nextval('spectrum_id'),
    session_name varchar(24) not null,
    session_index int not null,
    start_time timestamp default null,
    latitude float8 default 0,
    longitude float8 default 0,
    altitude float8 default 0,
    track float8 default 0,
    speed float8 default 0,
    climb float8 default 0,
    livetime float8 default 0,
    realtime float8 default 0,
    num_channels int default 0,
    channels text default '',
    doserate float8 default 0
);
alter sequence spectrum_id owned by spectrum.id;
alter table spectrum alter column start_time set default now();

grant all privileges on table spectrum to numsys;
grant usage, select on sequence spectrum_id to numsys;

