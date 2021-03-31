CREATE TABLE users (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  -- Stripe customer id
  stripe_id varchar(50),
  first_name varchar(40) NOT NULL,
  last_name varchar(40) NOT NULL,
  phone varchar(40) NOT NULL,
  email varchar(50) NOT NULL UNIQUE,
  verified boolean NOT NULL DEFAULT false,
  user_type varchar(4) NOT NULL CHECK (user_type IN ('GOAT', 'USER')),
  -- profile_avatar is a filepath.
  profile_avatar varchar(255),
  last_seen timestamp NOT NULL DEFAULT now(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now(),
  deleted_at timestamp
);

CREATE TABLE users_contacts (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  user_id int unsigned NOT NULL,
  contact_id int unsigned NOT NULL,
  request_status varchar(9) NOT NULL CHECK (request_status IN ('REQUESTED', 'ACCEPTED', 'DENIED')),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now(),

  UNIQUE KEY unique_users_contacts (user_id, contact_id),
  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (contact_id) REFERENCES users (id)
);

CREATE VIEW contacts AS
SELECT
	users_contacts.id,
  users_contacts.user_id,
  users_contacts.contact_id,
  users.first_name,
  users.last_name,
  users.phone,
  users.email,
  users.user_type,
  users.profile_avatar,
  users.last_seen,
  users_contacts.created_at,
  users_contacts.updated_at
FROM users_contacts JOIN users ON users_contacts.contact_id = users.id
WHERE request_status = 'ACCEPTED' AND users.deleted_at IS NULL;

CREATE VIEW contact_requests AS
SELECT
	users_contacts.id,
  users_contacts.user_id,
  users_contacts.contact_id,
  users_contacts.request_status,
  users.first_name,
  users.last_name,
  users.phone,
  users.email,
  users.user_type,
  users.profile_avatar,
  users.last_seen,
  users_contacts.created_at,
  users_contacts.updated_at
FROM users_contacts JOIN users ON users_contacts.contact_id = users.id
WHERE request_status IN ('REQUESTED', 'DENIED') AND users.deleted_at IS NULL;

CREATE TABLE credentials (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  user_id int unsigned NOT NULL UNIQUE,
  hash varchar(255) NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now(),

  FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE goat_invite_codes (
  -- Generated on row creation.
  invite_code char(6) PRIMARY KEY,
  -- User that requested invite code creation (admin or goat inviter).
  created_by int unsigned NOT NULL,
  status varchar(7) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'USED', 'EXPIRED')),
  -- The goat that registered using this code.
  used_by int unsigned UNIQUE,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now(),

  FOREIGN KEY (created_by) REFERENCES users (id),
  FOREIGN KEY (used_by) REFERENCES users (id)
);

CREATE TABLE posts (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  title varchar(255) NOT NULL,
  content text NOT NULL,
  -- Comma separated list of image urls.
  image_urls varchar(6553) NOT NULL DEFAULT '',
  tickers varchar(6553) NOT NULL DEFAULT '',
  posted_by int unsigned NOT NULL,
  published_at timestamp,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now(),

  FOREIGN KEY (posted_by) REFERENCES users (id)
);

CREATE TABLE comments (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  post_id int unsigned NOT NULL,
  content text NOT NULL,
  posted_by int unsigned NOT NULL,
  published_at timestamp,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now(),

  FOREIGN KEY (post_id) REFERENCES posts (id),
  FOREIGN KEY (posted_by) REFERENCES users (id)
);

CREATE TABLE threads (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  title varchar(255) DEFAULT '',
  thread_type varchar(6) NOT NULL CHECK (thread_type in ('SINGLE', 'GROUP')),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE TABLE thread_participants (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  user_id int unsigned NOT NULL REFERENCES users (id),
  thread_id int unsigned NOT NULL REFERENCES threads (id),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now(),

  FOREIGN KEY (user_id) REFERENCES users (id),
  FOREIGN KEY (thread_id) REFERENCES threads (id)
);

CREATE TABLE thread_messages (
  id int unsigned AUTO_INCREMENT PRIMARY KEY,
  thread_id int unsigned NOT NULL REFERENCES threads (id),
  sender_id int unsigned NOT NULL REFERENCES users (id),
  message_type varchar(4) NOT NULL DEFAULT 'TEXT' CHECK (message_type in ('TEXT')),
  -- pubnub timestamp.
  message text NOT NULL,
  timestamp bigint unsigned NOT NULL,

  FOREIGN KEY (thread_id) REFERENCES threads (id),
  FOREIGN KEY (sender_id) REFERENCES users (id)
);