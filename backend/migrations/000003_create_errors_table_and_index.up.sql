CREATE TABLE errors (
    site varchar(15),
    id integer,
    data text
);

CREATE unique index errors_index on errors(site, id);
