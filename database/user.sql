CREATE TABLE t_user (
    id serial primary key,
    uuid varchar(64) not null unique,
    email varchar(255) not null unique,
    name varchar(255),
    password varchar(255) not null,
    created_at timestamp not null,
    available boolean not null
);