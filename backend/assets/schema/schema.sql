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

create table errors (
    site varchar(15),
    id integer,
    data text
);

create unique index errors_index on errors(site, id);