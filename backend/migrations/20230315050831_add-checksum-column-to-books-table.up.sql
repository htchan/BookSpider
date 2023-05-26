-- Add checksum column to books table. The column will be bas64_encode(simplified_chinese("title"-"writer_name"))
ALTER TABLE public.books ADD COLUMN checksum TEXT;

-- Add a unique index to the checksum to ensure there are no duplicates
CREATE INDEX checksum_index ON public.books(checksum);
