-- Enable UUID extension (once per DB)
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Table: auth_users
CREATE TABLE auth_users (
  username TEXT PRIMARY KEY,
  full_name TEXT NOT NULL,
  email TEXT UNIQUE NOT NULL,
  hashed_password TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'user', 
  is_verified BOOLEAN NOT NULL DEFAULT FALSE,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Table: auth_session
CREATE TABLE auth_session (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username TEXT NOT NULL,
  refresh_token TEXT NOT NULL,
  user_agent TEXT NOT NULL,
  client_ip TEXT NOT NULL,
  is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
  expired_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT fk_auth_session_user FOREIGN KEY (username) REFERENCES auth_users (username) ON DELETE CASCADE
);

CREATE INDEX idx_auth_session_username ON auth_session (username);

-- Table: verify_email
CREATE TABLE verify_email (
  id BIGSERIAL PRIMARY KEY,
  username TEXT NOT NULL,
  email TEXT NOT NULL,
  secret_code TEXT NOT NULL,
  is_used BOOLEAN NOT NULL DEFAULT FALSE,
  expired_at TIMESTAMPTZ NOT NULL DEFAULT (now() + interval '15 minutes'),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT fk_verify_email_user FOREIGN KEY (username) REFERENCES auth_users (username) ON DELETE CASCADE
);

CREATE INDEX idx_verify_email_username ON verify_email (username);
