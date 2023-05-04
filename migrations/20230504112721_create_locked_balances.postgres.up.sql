CREATE TABLE locked_balances (
                                 id uuid NOT NULL,
                                 created_at bigint NOT NULL,
                                 updated_at bigint NOT NULL,
                                 deleted_at bigint DEFAULT NULL,
                                 account_id uuid NOT NULL,
                                 lock_date bigint NOT NULL,
                                 release_date bigint NOT NULL,
                                 amount_locked bigint NOT NULL,
                                 PRIMARY KEY (id),
                                 FOREIGN KEY (account_id) REFERENCES accounts(id)
);