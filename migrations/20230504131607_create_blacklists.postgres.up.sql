CREATE TABLE black_lists (
                             id uuid NOT NULL,
                             created_at bigint NOT NULL,
                             updated_at bigint NOT NULL,
                             deleted_at bigint DEFAULT NULL,
                             token varchar(255) NOT NULL,
                             email varchar(255) NOT NULL,
                             PRIMARY KEY (id)
);