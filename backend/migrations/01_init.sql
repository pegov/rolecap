CREATE TABLE IF NOT EXISTS auth_user (
	id SERIAL PRIMARY KEY,
	username TEXT,
	email TEXT,
	password TEXT,
	active BOOLEAN,
	verified BOOLEAN,

	created_at TIMESTAMP WITH TIME ZONE,
	last_login TIMESTAMP WITH TIME ZONE
);
CREATE INDEX IF NOT EXISTS auth_user_username_idx ON auth_user(username);
CREATE INDEX IF NOT EXISTS auth_user_email_idx ON auth_user(email);

CREATE TABLE IF NOT EXISTS auth_role (
	id SERIAL PRIMARY KEY,
	name TEXT
);
CREATE INDEX IF NOT EXISTS auth_role_name_idx ON auth_role(name);

CREATE TABLE IF NOT EXISTS auth_permission (
  id SERIAL PRIMARY KEY,
  name TEXT
);
CREATE INDEX IF NOT EXISTS auth_permission_name_idx ON auth_permission(name);

CREATE TABLE IF NOT EXISTS auth_user_role (
	user_id INTEGER REFERENCES auth_user(id) ON DELETE CASCADE,
	role_id INTEGER REFERENCES auth_role(id) ON DELETE CASCADE,
	PRIMARY KEY(user_id, role_id)
);
CREATE INDEX IF NOT EXISTS auth_user_role_user_id_idx ON auth_user_role(user_id);
CREATE INDEX IF NOT EXISTS auth_user_role_role_id_idx ON auth_user_role(role_id);

CREATE TABLE IF NOT EXISTS auth_role_permission (
  role_id INTEGER REFERENCES auth_role(id) ON DELETE CASCADE,
  permission_id INTEGER REFERENCES auth_permission(id) ON DELETE CASCADE,
	PRIMARY KEY(role_id, permission_id)
);
CREATE INDEX IF NOT EXISTS auth_role_permission_role_id_idx ON auth_role_permission(role_id);
CREATE INDEX IF NOT EXISTS auth_role_permission_permission_id_idx ON auth_role_permission(permission_id);

CREATE TABLE IF NOT EXISTS auth_oauth (
	user_id INTEGER REFERENCES auth_user(id) ON DELETE CASCADE,
	provider TEXT,
	sid TEXT,
  PRIMARY KEY(user_id)
);
CREATE INDEX IF NOT EXISTS auth_oauth_provider_sid_idx ON auth_oauth(provider, sid);