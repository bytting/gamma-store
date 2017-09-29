drop database gs;
create database gs;

\c gs;

create sequence spectrum_id;

create table spectrum (
    id int not null default nextval('spectrum_id'),
    session_name varchar(24) not null,
    session_index int not null,
    start_time varchar(24) default '',
    latitude float8 default 0,
    latitude_error float8 default 0,
    longitude float8 default 0,
    longitude_error float8 default 0,
    altitude float8 default 0,
    altitude_error float8 default 0,
    track float8 default 0,
    track_error float8 default 0,
    speed float8 default 0,
    speed_error float8 default 0,
    climb float8 default 0,
    climb_error float8 default 0,
    livetime float8 default 0,
    realtime float8 default 0,
    total_count int default 0,
    num_channels int default 0,
    channels text default '',
    doserate float8 default 0
);

alter sequence spectrum_id owned by spectrum.id;

grant all privileges on table spectrum to numsys;
grant usage, select on sequence spectrum_id to numsys;
