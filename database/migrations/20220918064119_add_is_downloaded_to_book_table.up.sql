ALTER TABLE books ADD COLUMN IF NOT EXISTS is_downloaded BOOLEAN DEFAULT false;
UPDATE books SET 
  is_downloaded=true, "status"='END'
WHERE "status"='DOWNLOAD';
