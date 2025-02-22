DROP TRIGGER IF EXISTS update_book_detail_timestamp ON book_details;
DROP TRIGGER IF EXISTS book_detail_tsvector_update ON book_details;

DROP INDEX IF EXISTS idx_book_details_available_copies;
DROP INDEX IF EXISTS idx_book_details_title_tsvector;
DROP INDEX IF EXISTS idx_book_details_title_lower_trimmed;

DROP TABLE IF EXISTS book_details;
