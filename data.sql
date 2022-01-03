
insert into genres (name, active)
values ('Novel', default),
       ('Novella', default),
       ('Tale', default),
       ('Fable', default),
       ('Fairy tale', default),
       ('Detective', default),
       ('Science fiction', default),
       ('Non-fiction', default),
       ('Mythology', default),
       ('Poem', default),
       ('Biography', default),
       ('Manual', default),
       ('Historical', default),
       ('Note', default),
       ('Fantasy', default);

insert into users (name, login, password, active, created)
values ('Tester', 'writer', 'password', default, default);

insert into books (title, author_id, genre_id, description, cover_image_name, access_read, active, created)
values ('Test Book', 1, 5, 'The first book', '6d95e2e7-ccbc-4bf7-828a-927f88225857.png', default, default, default);

insert into chapters (book_id, number, name, content, active, created)
values (1, 1, 'Beginning', 'Once upon a time...', default, default);

insert into chapters (book_id, number, name, content, active, created)
values  (1, 2, 'End', '...and they lived happily ever after. ', default, default)
