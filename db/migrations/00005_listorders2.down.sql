BEGIN TRANSACTION;
ALTER TABLE orders
    DROP COLUMN IF EXISTS accrual,
    DROP COLUMN IF EXISTS uploaded_at;
COMMIT;