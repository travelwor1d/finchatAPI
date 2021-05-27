-- +goose Up
-- SQL in this section is executed when the migration is applied.


DROP VIEW IF EXISTS verified_active_users;



alter table core.users drop column updated_at;


alter table core.users add column updated_at timestamp NOT NULL DEFAULT now();

drop TRIGGER users_updated_at;
CREATE TRIGGER users_updated_at BEFORE UPDATE ON core.users
FOR EACH ROW
BEGIN
    IF NEW.last_seen = OLD.last_seen THEN
      SET  NEW.updated_at  = now();
    END IF;
END


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.