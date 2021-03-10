CREATE TABLE users (
  id int AUTO_INCREMENT PRIMARY KEY,
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
  updated_at timestamp NOT NULL DEFAULT now() ON UPDATE now()
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
  used_by int UNIQUE REFERENCES users (id)
);