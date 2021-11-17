create table users
(
    id       bigserial primary key,
    name     text      not null,
    login    text      not null unique,
    password text      not null unique,
    active   boolean   not null default true,
    created  timestamp not null default current_timestamp
);

create table books
(
    id          bigserial primary key,
    title       text      not null,
    author_id   bigint    not null references users,
    description text,
    cover_image text      not null,
    access_read boolean not null default true,
    active boolean not null  default true,
    created     timestamp not null default current_timestamp
)