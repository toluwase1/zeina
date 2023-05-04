CREATE TABLE transactions (
                              id uuid NOT NULL,
                              created_at bigint NOT NULL,
                              updated_at bigint NOT NULL,
                              deleted_at bigint DEFAULT NULL,
                              account_id uuid NOT NULL,
                              entry varchar(255) NOT NULL,
                              purpose varchar(255) NOT NULL,
                              description varchar(255) NOT NULL,
                              remark varchar(255) NOT NULL,
                              status varchar(255) NOT NULL,
                              beneficiary_name varchar(255) NOT NULL,
                              total_balance bigint NOT NULL,
                              available_balance bigint NOT NULL,
                              pending_balance bigint NOT NULL,
                              PRIMARY KEY (id),
                              FOREIGN KEY (account_id) REFERENCES accounts(id)
);