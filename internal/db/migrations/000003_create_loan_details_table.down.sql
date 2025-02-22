DROP TRIGGER IF EXISTS update_loan_details_timestamp ON loan_details;

DROP INDEX IF EXISTS idx_loan_details_loan_date;
DROP INDEX IF EXISTS idx_loan_details_return_date;
DROP INDEX IF EXISTS idx_loan_details_book_id;

ALTER TABLE loan_details
  DROP CONSTRAINT IF EXISTS fk_loan_details_book;

DROP TABLE IF EXISTS loan_details;
