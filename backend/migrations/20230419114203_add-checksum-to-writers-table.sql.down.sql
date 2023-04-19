-- Remove the unique index on the checksum column
DROP INDEX IF EXISTS writers_checksum_index;

-- Remove checksum column from writers table
ALTER TABLE public.writers DROP COLUMN checksum;

-- Add writer checksum column to books table.
ALTER TABLE public.books DROP COLUMN writer_checksum TEXT;
DROP index books_writer_checksum on books(writer_checksum);
