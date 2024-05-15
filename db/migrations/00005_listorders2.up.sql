BEGIN TRANSACTION;
ALTER TABLE orders
    ADD COLUMN accrual     int                      DEFAULT 0,
    ADD COLUMN uploaded_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP;
COMMIT;