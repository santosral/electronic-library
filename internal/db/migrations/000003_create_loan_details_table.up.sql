CREATE TABLE loan_details (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name_of_borrower VARCHAR(255) NOT NULL,
    loan_date DATE,
    return_date DATE,
    book_id UUID NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_loan_details_book FOREIGN KEY (book_id) REFERENCES book_details(id) 
    ON DELETE CASCADE
    ON UPDATE CASCADE
);

CREATE INDEX idx_loan_details_loan_date ON loan_details(loan_date);
CREATE INDEX idx_loan_details_return_date ON loan_details(return_date);
CREATE INDEX idx_loan_details_book_id ON loan_details(book_id);

CREATE TRIGGER update_book_detail_timestamp
BEFORE UPDATE ON loan_details
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
