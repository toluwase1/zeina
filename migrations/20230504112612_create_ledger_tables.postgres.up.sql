CREATE TABLE ledgers (
                         id uuid NOT NULL,
                         created_at bigint NOT NULL,
                         account_id uuid NOT NULL,
                         account_type varchar(255) NOT NULL,
                         entry varchar(255) NOT NULL,
                         change bigint NOT NULL,
                         type varchar(255) NOT NULL,
                         PRIMARY KEY (id),
                         FOREIGN KEY (account_id) REFERENCES accounts(id)
);