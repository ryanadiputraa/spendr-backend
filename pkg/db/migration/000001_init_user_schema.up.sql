CREATE TABLE users (
  id VARCHAR(256) PRIMARY KEY,
  email VARCHAR(256) UNIQUE NOT NULL,
  password VARCHAR(256),
  first_name VARCHAR(50) NOT NULL,
  last_name VARCHAR(50) NOT NULL,
  picture VARCHAR(256),
  currency VARCHAR(10),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users("email");
