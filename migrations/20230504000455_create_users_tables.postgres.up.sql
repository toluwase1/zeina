CREATE TABLE users (
                       id uuid NOT NULL,
                       email varchar(255) NOT NULL,
                       name varchar(255) NOT NULL,
                       created_at bigint NOT NULL,
                       updated_at bigint NOT NULL,
                       deleted_at bigint DEFAULT NULL,
                       PRIMARY KEY (id)
);
