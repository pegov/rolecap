-- name: get_by_id
SELECT * FROM auth_user WHERE id = $1;

-- name: get_by_email
SELECT * FROM auth_user WHERE email = $1;

-- name: get_by_username
SELECT * FROM auth_user WHERE username = $1;

-- name: create
INSERT INTO auth_user(
	email,
	username,
	password,
	active,
	verified,
	created_at,
	last_login
) VALUES (
	$1, $2, $3, $4, $5, $6, $7
) RETURNING id;