-- +goose Up
-- SQL in this section is executed when the migration is applied.



alter table core.users_contacts add (uuid VARCHAR(36));



-- +goose Down
-- SQL in this section is executed when the migration is rolled back.