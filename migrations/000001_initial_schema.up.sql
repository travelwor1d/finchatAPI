CREATE TABLE users (
  id int AUTO_INCREMENT PRIMARY KEY,
  first_name varchar(40) NOT NULL,
  last_name varchar(40) NOT NULL,
  phone varchar(40),
  email varchar(50) NOT NULL,
  user_type varchar(4) NOT NULL CHECK (user_type in ('GOAT', 'USER')),
  -- profile_avatar is a filepath.
  profile_avatar varchar(255),
  last_seen timestamp NOT NULL DEFAULT now()
);