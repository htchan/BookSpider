CREATE TABLE writers (
    -- id integer primary key AUTOINCREMENT,
    id serial primary key,
    name text
);

CREATE unique index writers_id on writers(id);
CREATE unique index writers_name on writers(name);
insert into writers (id, name) values (0, '');
