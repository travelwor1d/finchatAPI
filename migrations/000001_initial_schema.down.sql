DROP VIEW IF EXISTS posts_with_votes;
DROP VIEW IF EXISTS active_users, verified_active_users, contacts;

DROP TABLE IF EXISTS thread_messages, thread_participants; -- Probably tables are not needed.
DROP TABLE IF EXISTS threads; -- Probably table is not needed.
DROP TABLE IF EXISTS comments, post_votes;
DROP TABLE IF EXISTS posts, goat_invite_codes, users_contacts;
DROP TABLE IF EXISTS users;