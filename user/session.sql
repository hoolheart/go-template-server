/** user visit session table */
create table t_session (
  id serial primary key, /**< record ID */
  uuid varchar(64) not null unique, /**< session UUID */
  user_id varchar(64) not null, /**< relevant user UUID */
  created_at timestamp not null /**< created timestamp */
);
