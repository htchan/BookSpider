-- Add checksum column to writers table. The column will be bas64_encode(simplified_chinese("writer_name"))
ALTER TABLE public.writers ADD COLUMN checksum TEXT;

-- Add a unique index to the checksum to ensure there are no duplicates
CREATE INDEX writers_checksum_index ON public.writers(checksum);

-- Add writer checksum column to books table.
ALTER TABLE public.books ADD COLUMN writer_checksum TEXT;
CREATE index books_writer_checksum on books(writer_checksum);