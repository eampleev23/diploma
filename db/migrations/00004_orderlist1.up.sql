BEGIN TRANSACTION;
ALTER TABLE orders
    ADD COLUMN status char(20);
COMMIT;