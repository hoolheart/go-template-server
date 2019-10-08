/** user event record table */
CREATE TABLE t_user_event (
    id serial primary key, /**< record ID */
    ts timestamp not null, /**< event timestamp */
    user_id varchar(64) not null, /**< relevent user UUID */
    /** event type: \
        * 0 --- none/invalid \
        * 1 --- register \
        * 2 --- login \
        * 3 --- logout \
        * 4 --- modify \
        * 5 --- destroy
     */
    type smallint not null,
    info text, /**< event information (meta-data in json) */
    login_ip varchar(50), /**< login IP */
    browser varchar(255) /**< used browser */
);