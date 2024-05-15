BEGIN TRANSACTION;
ALTER TABLE orders
    ALTER COLUMN number TYPE INT8;
CREATE UNIQUE INDEX IF NOT EXISTS order_number_unique
    ON orders
        USING btree (number);
ALTER TABLE orders
    ADD COLUMN customer_id INT;
ALTER TABLE orders
    ADD FOREIGN KEY (customer_id)
        REFERENCES users (id);
COMMIT;