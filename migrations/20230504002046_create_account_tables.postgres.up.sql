CREATE TABLE accounts (
                          id uuid NOT NULL,
                          created_at bigint NOT NULL,
                          updated_at bigint NOT NULL,
                          deleted_at bigint DEFAULT NULL,
                          user_id uuid NOT NULL,
                          active boolean NOT NULL,
                          total_balance bigint NOT NULL,
                          available_balance bigint NOT NULL,
                          pending_balance bigint NOT NULL,
                          locked_balance bigint NOT NULL,
                          PRIMARY KEY (id),
                          FOREIGN KEY (user_id) REFERENCES users(id)
);