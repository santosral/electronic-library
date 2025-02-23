CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE book_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    available_copies INT NOT NULL DEFAULT 0,
    title_tsvector TSVECTOR,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_book_details_available_copies ON book_details(available_copies);
CREATE INDEX idx_book_details_title_tsvector ON book_details USING GIN(to_tsvector('english', title));
CREATE UNIQUE INDEX idx_book_details_title_lower_trimmed ON book_details (LOWER(TRIM(title)));

CREATE TRIGGER update_book_detail_timestamp
BEFORE UPDATE ON book_details
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TRIGGER book_detail_tsvector_update
BEFORE INSERT OR UPDATE ON book_details
FOR EACH ROW
EXECUTE FUNCTION tsvector_update_trigger(title_tsvector, 'pg_catalog.english', title);