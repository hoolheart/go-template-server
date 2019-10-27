/** user table */
CREATE TABLE t_user (
  id serial primary key, /**< record ID */
  uuid varchar(64) not null unique, /**< user UUID */
  email varchar(255) not null unique, /**< user email */
  name varchar(255), /**< user name */
  password varchar(255) not null, /**< user password (encrupted) */
  created_at timestamp not null /**< created timestamp */
);