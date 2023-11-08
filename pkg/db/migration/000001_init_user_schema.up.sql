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

CREATE TABLE expense_categories (
  id VARCHAR(256) PRIMARY KEY,
  category VARCHAR(100) UNIQUE NOT NULL,
  ico VARCHAR(256) NOT NULL,
  user_id VARCHAR(256) NOT NULL,
  CONSTRAINT fk_user_id
    FOREIGN KEY(user_id)
      REFERENCES users(id)
);

CREATE TABLE expenses (
  id VARCHAR(256) PRIMARY KEY,
  category_id VARCHAR(256),
  user_id VARCHAR(256) NOT NULL,
  expense VARCHAR(100) NOT NULL,
  amount INT NOT NULL,
  date TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT fk_exepnse_category
    FOREIGN KEY(category_id)
      REFERENCES expense_categories(id),
  CONSTRAINT fk_user_id
    FOREIGN KEY(user_id)
      REFERENCES users(id)
);

CREATE INDEX idx_users_email ON users("email");
CREATE INDEX idx_expenses_date ON expenses("date");
