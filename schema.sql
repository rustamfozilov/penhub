create table users
(
    id       bigserial primary key,
    name     text      not null,
    login    text      not null unique,
    password text      not null unique,
    active   boolean   not null default true,
    created  timestamp not null default current_timestamp
);

create table users_tokens
(
    user_id bigint      not null references users,
    token   text        not null unique,
    expire  timestamptz not null default current_timestamp + interval '1 hour',
    created timestamptz not null default current_timestamp
);

create table books
(
    id          bigserial primary key,
    title       text      not null,
    author_id   bigint    not null references users,
    genre_id       bigint      not null references genres,
    description text not null default 'description',
    cover_image text      not null,
    access_read boolean   not null default true,
    active      boolean   not null default true,
    created     timestamp not null default current_timestamp
);
drop table books cascade;

create table chapters
(
    id      bigserial primary key,
    book_id bigint    not null references books,
    number  bigint    not null,
    name    text      not null,
    content text      not null,
    active  boolean   not null default true,
    created timestamp not null default current_timestamp
);

create table genres
(
    id     bigserial primary key,
    name   text    not null,
    active boolean not null default true
);