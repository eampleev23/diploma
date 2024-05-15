BEGIN TRANSACTION;
ALTER TABLE withdraw
    ADD COLUMN user_id INT;
ALTER TABLE withdraw
    ADD FOREIGN KEY (user_id)
        REFERENCES users (id);
COMMIT;