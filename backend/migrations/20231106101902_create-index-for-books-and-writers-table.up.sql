CREATE INDEX IF NOT EXISTS books__vendor_reference ON books(site,id,hash_code desc);
CREATE INDEX IF NOT EXISTS books__status ON books(status,is_downloaded);
CREATE INDEX IF NOT EXISTS books__checksum ON books (checksum, writer_checksum);

CREATE INDEX IF NOT EXISTS writers__name ON writers(name);
