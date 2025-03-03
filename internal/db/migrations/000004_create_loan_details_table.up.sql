CREATE TABLE LOAN_DETAILS (
	ID UUID PRIMARY KEY DEFAULT UUID_GENERATE_V4 (),
	NAME_OF_BORROWER VARCHAR(255) NOT NULL,
	LOAN_DATE DATE NOT NULL,
	RETURN_DATE DATE NOT NULL,
	RETURNED_ON DATE,
	BOOK_ID UUID NOT NULL,
	CREATED_AT TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	UPDATED_AT TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT FK_LOAN_DETAILS_BOOK FOREIGN KEY (BOOK_ID) REFERENCES BOOK_DETAILS (ID) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE INDEX IDX_LOAN_DETAILS_LOAN_DATE ON LOAN_DETAILS (LOAN_DATE);

CREATE INDEX IDX_LOAN_DETAILS_RETURN_DATE ON LOAN_DETAILS (RETURN_DATE);

CREATE INDEX IDX_LOAN_DETAILS_RETURNED_ON ON LOAN_DETAILS (RETURNED_ON);

CREATE INDEX IDX_LOAN_DETAILS_BOOK_ID ON LOAN_DETAILS (BOOK_ID);

CREATE TRIGGER UPDATE_BOOK_DETAIL_TIMESTAMP BEFORE
UPDATE ON LOAN_DETAILS FOR EACH ROW
EXECUTE FUNCTION UPDATE_TIMESTAMP ();