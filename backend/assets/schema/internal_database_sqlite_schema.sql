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
(1, 'writer-1'),
(2, 'writer-2'),
(3, 'writer-3');

insert into errors
(site, id, data)
values
('test', 2, 'error-2'),
('test', 4, 'error-4'),
('test', 5, 'error-5');

