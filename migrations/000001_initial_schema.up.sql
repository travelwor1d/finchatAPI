CREATE TABLE users (
  id int AUTO_INCREMENT PRIMARY KEY,
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

CREATE TABLE credentials (
  id int AUTO_INCREMENT PRIMARY KEY,
  user_id int NOT NULL UNIQUE REFERENCES users (id),
  hash varchar(255) NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE TABLE goat_invite_codes (
  -- Generated on row creation.
  invite_code char(6) PRIMARY KEY,
  -- User that requested invite code creation (admin or goat inviter).
  created_by int NOT NULL REFERENCES users (id),
  status varchar(7) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'USED', 'EXPIRED')),
  -- The goat that registered using this code.
  used_by int UNIQUE REFERENCES users (id),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE TABLE posts (
  id int AUTO_INCREMENT PRIMARY KEY,
  title varchar(255) NOT NULL,
  content text NOT NULL,
  -- Comma separeted list of image urls.
  image_urls varchar(6553) NOT NULL DEFAULT '',
  tickers varchar(6553) NOT NULL DEFAULT '',
  posted_by int NOT NULL REFERENCES users (id),
  published_at timestamp,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE TABLE comments (
  id int AUTO_INCREMENT PRIMARY KEY,
  post_id int NOT NULL REFERENCES posts (id),
  content text NOT NULL,
  posted_by int NOT NULL REFERENCES users (id),
  published_at timestamp,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE TABLE threads (
  id int AUTO_INCREMENT PRIMARY KEY,
  title varchar(255) DEFAULT '',
  thread_type varchar(6) NOT NULL CHECK (thread_type in ('SINGLE', 'GROUP')),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE TABLE thread_participants (
  id int AUTO_INCREMENT PRIMARY KEY,
  user_id int NOT NULL REFERENCES users (id),
  thread_id int NOT NULL REFERENCES threads (id),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE TABLE thread_messages (
  id int AUTO_INCREMENT PRIMARY KEY,
  thread_id int NOT NULL REFERENCES threads (id),
  sender_id int NOT NULL REFERENCES users (id),
  message_type varchar(4) NOT NULL DEFAULT 'TEXT' CHECK (message_type in ('TEXT')),
  -- pubnub timestamp.
  message text NOT NULL,
  timestamp bigint NOT NULL
);