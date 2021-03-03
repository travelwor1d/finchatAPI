CREATE TABLE users (
  id int AUTO_INCREMENT PRIMARY KEY,
  first_name varchar(40) NOT NULL,
  last_name varchar(40) NOT NULL,
  phone varchar(40),
  email varchar(50) NOT NULL UNIQUE,
  user_type varchar(4) NOT NULL CHECK (user_type in ('GOAT', 'USER')),
  -- profile_avatar is a filepath.
  profile_avatar varchar(255),
  last_seen timestamp NOT NULL DEFAULT now(),
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now()
);

CREATE TABLE credentials (
  id int AUTO_INCREMENT PRIMARY KEY,
  user_id int NOT NULL UNIQUE REFERENCES users (id),
  hash varchar(255) NOT NULL,
  created_at timestamp NOT NULL DEFAULT now(),
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now()
);