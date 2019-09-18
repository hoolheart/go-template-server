CREATE TABLE t_user_event (
    id serial primary key,
    ts timestamp not null,
    user_id varchar(64) not null,
    type smallint not null,
    info text,
    login_ip varchar(50),
    browser varchar(255)
);