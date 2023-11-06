UPDATE books SET books.status='DOWNLOAD' WHERE books.is_downloaded=true;
ALTER TABLE books DROP COLUMN IF EXISTS is_downloaded;