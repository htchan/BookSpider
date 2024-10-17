-- name: CreateBookWithZeroHash :one
INSERT INTO books
(site, id, hash_code, title, writer_id, writer_checksum, type, 
update_date, update_chapter, status, is_downloaded, checksum)
VALUES
($1, $2, 0, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING *;

-- name: CreateBookWithHash :one
INSERT INTO books
(site, id, hash_code, title, writer_id, writer_checksum, type, 
update_date, update_chapter, status, is_downloaded, checksum)
VALUES
($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: UpdateBook :one
Update books SET 
title=$4, writer_id=$5, writer_checksum=$12, type=$6, update_date=$7, update_chapter=$8,
status=$9, is_downloaded=$10, checksum=$11
WHERE site=$1 and id=$2 and hash_code=$3
RETURNING *;

-- name: GetBookByID :one
select books.site, books.id, books.hash_code, books.title,
  books.writer_id, coalesce(writers.name, ''), books.type,
  books.update_date, books.update_chapter, 
  books.status, books.is_downloaded, coalesce(errors.data, '')
from books left join writers on books.writer_id=writers.id 
  left join errors on books.site=errors.site and books.id=errors.id
where books.site=$1 and books.id=$2 order by books.hash_code desc;

-- name: GetBookByIDHash :one
select books.site, books.id, books.hash_code, books.title,
  books.writer_id, coalesce(writers.name, ''), books.type,
  books.update_date, books.update_chapter, 
  books.status, books.is_downloaded, coalesce(errors.data, '')
from books left join writers on books.writer_id=writers.id 
  left join errors on books.site=errors.site and books.id=errors.id
where books.site=$1 and books.id=$2 and books.hash_code=$3
order by hash_code desc;

-- name: ListBooksByStatus :many
select books.site, books.id, books.hash_code, books.title,
  books.writer_id, coalesce(writers.name, ''), books.type,
  books.update_date, books.update_chapter, 
  books.status, books.is_downloaded, coalesce(errors.data, '')
from books left join writers on books.writer_id=writers.id 
  left join errors on books.site=errors.site and books.id=errors.id
where books.site=$1 and books.status=$2 order by hash_code desc;

-- name: ListBooks :many
select books.site, books.id, books.hash_code, books.title,
  books.writer_id, coalesce(writers.name, ''), books.type,
  books.update_date, books.update_chapter, 
  books.status, books.is_downloaded, coalesce(errors.data, '')
from books left join writers on books.writer_id=writers.id 
  left join errors on books.site=errors.site and books.id=errors.id
where books.site=$1
order by books.site, books.id, books.hash_code;

-- name: ListBooksForUpdate :many
select distinct on (books.site, books.id) 
  books.site, books.id, books.hash_code, books.title,
  books.writer_id, coalesce(writers.name, ''), books.type,
  books.update_date, books.update_chapter, 
  books.status, books.is_downloaded, coalesce(errors.data, '')
from books left join writers on books.writer_id=writers.id 
  left join errors on books.site=errors.site and books.id=errors.id
where books.site=$1
order by books.site, books.id desc, books.hash_code desc;

-- name: ListBooksForDownload :many
select distinct on (books.site, books.id) 
  books.site, books.id, books.hash_code, books.title,
  books.writer_id, coalesce(writers.name, ''), books.type,
  books.update_date, books.update_chapter, 
  books.status, books.is_downloaded, coalesce(errors.data, '')
from books left join writers on books.writer_id=writers.id 
  left join errors on books.site=errors.site and books.id=errors.id
where books.site=$1 and books.status='END' and books.is_downloaded=false
order by books.site, books.id desc, books.hash_code desc;

-- name: ListBooksByTitleWriter :many
select books.site, books.id, books.hash_code, books.title,
  books.writer_id, coalesce(writers.name, ''), books.type,
  books.update_date, books.update_chapter, 
  books.status, books.is_downloaded, coalesce(errors.data, '')
from books left join writers on books.writer_id=writers.id
  left join errors on books.site=errors.site and books.id=errors.id
where books.site=$1 and books.status != 'ERROR' and 
  (($2 != '%%' and books.title like $2) or
  ($3 != '%%' and writers.name like $3))
order by books.update_date desc, books.id desc limit $4 offset $5;

-- name: ListRandomBooks :many
select books.site, books.id, books.hash_code, books.title,
  books.writer_id, coalesce(writers.name, ''), books.type,
  books.update_date, books.update_chapter, 
  books.status, books.is_downloaded, coalesce(errors.data, '')
from books left join writers on books.writer_id=writers.id 
  left join errors on books.site=errors.site and books.id=errors.id
where books.site=$1 and books.is_downloaded=true
order by books.site, books.id desc, books.hash_code desc 
limit $2 offset RANDOM() * 
(
  select greatest(count(*) - $2, 0)
  from books as bks where site=$1 and bks.is_downloaded=true
);

-- name: GetBookGroupByID :many
select books.site, books.id, books.hash_code, books.title,
  books.writer_id, coalesce(writers.name, ''), books.type,
  books.update_date, books.update_chapter, 
  books.status, books.is_downloaded, coalesce(errors.data, '')
from books
  left join writers on books.writer_id=writers.id 
  left join errors on books.site=errors.site and books.id=errors.id
where (books.checksum, books.writer_checksum) = (
  select bks.checksum, bks.writer_checksum from books as bks 
  where bks.site=$1 and bks.id=$2 
  and bks.checksum != '' and bks.writer_checksum != ''
  order by bks.hash_code desc limit 1
) or books.site=$1 and books.id=$2;

-- name: GetBookGroupByIDHash :many
select books.site, books.id, books.hash_code, books.title,
  books.writer_id, coalesce(writers.name, ''), books.type,
  books.update_date, books.update_chapter, 
  books.status, books.is_downloaded, coalesce(errors.data, '')
from books
  left join writers on books.writer_id=writers.id 
  left join errors on books.site=errors.site and books.id=errors.id
where (books.checksum, books.writer_checksum) = (
  select bks.checksum, bks.writer_checksum from books as bks 
  where bks.site=$1 and bks.id=$2 and bks.hash_code=$3 
  and bks.checksum != '' and bks.writer_checksum != ''
  order by bks.hash_code desc limit 1
) or books.site=$1 and books.id=$2;

-- name: UpdateBooksStatus :exec
update books set is_downloaded=false, status='END' 
where (update_date < $1 or 
    update_chapter like '%番外%' or update_chapter like '%結局%' or 
    update_chapter like '%新書%' or update_chapter like '%完結%' or 
    update_chapter like '%尾聲%' or update_chapter like '%感言%' or 
    update_chapter like '%後記%' or update_chapter like '%完本%' or 
    update_chapter like '%全書完%' or update_chapter like '%全文完%' or 
    update_chapter like '%全文終%' or update_chapter like '%全文結%' or 
    update_chapter like '%劇終%' or update_chapter like '%（完）%' or 
    update_chapter like '%終章%' or update_chapter like '%外傳%' or 
    update_chapter like '%結尾%' or update_chapter like '%番外%' or 
    update_chapter like '%结局%' or update_chapter like '%新书%' or 
    update_chapter like '%完结%' or update_chapter like '%尾声%' or 
    update_chapter like '%感言%' or update_chapter like '%后记%' or 
    update_chapter like '%完本%' or update_chapter like '%全书完%' or 
    update_chapter like '%全文完%' or update_chapter like '%全文终%' or 
    update_chapter like '%全文结%' or update_chapter like '%剧终%' or 
    update_chapter like '%（完）%' or update_chapter like '%终章%' or 
    update_chapter like '%外传%' or update_chapter like '%结尾%') and 
  status='INPROGRESS' and site=$2;

-- name: CreateWriter :one
insert into writers (name, checksum) values ($1, $2) 
on conflict (name) do update set name=$1 
returning *;

-- name: CreateError :one
insert into errors (site, id, data) values ($1, $2, $3)
on conflict (site, id)
do update set data=$3
RETURNING *;

-- name: DeleteError :one
delete from errors where site=$1 and id=$2 returning *;

-- name: BackupBooks :exec
copy (select * from books where site=$1) to '$2' 
csv header quote as '''' force quote *;

-- name: BackupWriters :exec
copy (
  select distinct(writers.*) from writers join books on writers.id=books.writer_id 
  where books.site=$1
) to '$2' csv header quote as '''' force quote *;

-- name: BackupError :exec
copy (select * from errors where site=$1) to '$2' 
csv header quote as '''' force quote *;

-- name: BooksStat :one
select count(*) as book_count, count(distinct id) as unique_book_count, max(id) as max_book_id from books where site=$1;

-- name: NonErrorBooksStat :one
select max(id) as latest_success_id from books where status<>'ERROR' and site=$1;

-- name: ErrorBooksStat :one
select count(*) as error_count from books where site=$1 and status='ERROR';

-- name: DownloadedBooksStat :one
select count(*) as downloaded_count from books where site=$1 and is_downloaded=true;

-- name: BooksStatusStat :many
select status, count(*) from books where site=$1 group by status;

-- name: WritersStat :one
select count(distinct writers_id) as writer_count 
from books where site=$1;

-- name: FindAllBookIDs :many
select distinct(id) as book_id from books where site=$1 order by book_id;
