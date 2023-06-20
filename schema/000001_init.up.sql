CREATE TABLE users
(
    id serial not null unique,
    username varchar(255) not null unique,
    password_hash varchar(255) not null
);

CREATE TABLE sessions
(
    user_id     int references users(id) on delete cascade not null,
    token       varchar(255),
    expire_date timestamp
);

CREATE TABLE audit
(
    user_id    int references users (id) on delete cascade not null,
    event_date date,
    event      int
)