-- Remove the unique index on the checksum column
DROP INDEX IF EXISTS checksum_index;

-- Remove checksum column from books table
ALTER TABLE public.books DROP COLUMN checksum;
