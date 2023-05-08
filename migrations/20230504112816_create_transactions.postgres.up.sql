CREATE TABLE transactions (
                              id uuid NOT NULL,
                              created_at bigint NOT NULL,
                              updated_at bigint NOT NULL,
                              deleted_at bigint DEFAULT NULL,
                              account_id uuid NOT NULL,
                              entry varchar(255) NOT NULL,
                              purpose varchar(255) NOT NULL,
                              status varchar(255) NOT NULL,
                              change bigint DEFAULT NULL,
                              available_balance bigint NOT NULL,
                              pending_balance bigint NOT NULL,
                              reference varchar(255) NOT NULL,
                              PRIMARY KEY (id),
                              FOREIGN KEY (account_id) REFERENCES accounts(id),
                              CONSTRAINT unique_reference UNIQUE (reference)
);
