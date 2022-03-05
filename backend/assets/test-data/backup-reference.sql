create table books (
    site varchar(15),
    id integer,
    hash_code integer,
    title text,
    writer_id integer,
    type varchar(20),
    update_date varchar(30),
    update_chapter text,
    status integer
);

create unique index books_index on books(site, id, hash_code);
create index books_title on books(title);
create index books_writer on books(writer_id);
create index books_status on books(status);

create table writers (
    id integer primary key AUTOINCREMENT,
    name text
);

create unique index writers_id on writers(id);
create unique index writers_name on writers(name);
insert into writers (id, name) values (0, '');

create table errors (
    site varchar(15),
    id integer,
    data text
);

create unique index errors_index on errors(site, id);

insert into books
(site, id, hash_code, title, writer_id, type, update_date, update_chapter, status)
values
('test', 1, 100, 'title-1', 1, 'type-1', '104', 'chapter-1', 1),
('test', 2, 101, '', 0, '', '', '', 0),
('test', 3, 102, 'title-3', 2, 'type-3', '102', 'chapter-3', 3),
('test', 4, 101, '', 0, '', '', '', 0),
('test', 5, 101, '', 0, '', '', '', 0),
('test', 3, 200, 'title-3-new', 3, 'type-3-new', '100', 'chapter-3-new', 2);

insert into writers
(id, name)
values
(0, ''),
(1, 'writer-1'),
(2, 'writer-2'),
(3, 'writer-3');

insert into errors
(site, id, data)
values
('test', 2, 'error-2'),
('test', 4, 'error-4'),
('test', 5, 'error-5');

