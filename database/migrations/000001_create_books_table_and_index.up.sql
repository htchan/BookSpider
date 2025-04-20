CREATE TABLE books (
    site varchar(15),
    id integer,
    hash_code integer,
    title text,
    writer_id integer,
    type varchar(20),
    update_date varchar(30),
    update_chapter text,
    status varchar(10)
);

CREATE unique index books_index on books(site, id, hash_code);
CREATE index books_title on books(title);
CREATE index books_writer on books(writer_id);
CREATE index books_status on books(status);
