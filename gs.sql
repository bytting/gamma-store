drop database gs;
create database gs;

\c gs

create sequence session_id;

create table session (
    id int not null default nextval('session_id'),
    name char(15),
    comment varchar(256),
    detector_data text
);

alter sequence session_id owned by session.id;

create sequence spectrum_id;

create table spectrum (
    id int not null default nextval('spectrum_id'),
    session_id int not null,
    session_index int not null,
    start_time timestamp,
    latitude float8,
    latitude_error float8,
    longitude float8,
    longitude_error float8,
    altitude float8,
    altitude_error float8,
    track float8,
    track_error float8,
    speed float8,
    speed_error float8,
    climb float8,
    climb_error float8,
    livetime float8,
    realtime float8,
    total_count int,
    num_channels int,
    channels text,
    doserate float8
);

alter sequence spectrum_id owned by spectrum.id;
