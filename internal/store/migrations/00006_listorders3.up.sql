BEGIN TRANSACTION;
ALTER TABLE orders
    ALTER COLUMN number TYPE text;
COMMIT;