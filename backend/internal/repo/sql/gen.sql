-- name: get_by_id
SELECT * FROM gen WHERE id = $1;

-- name: create
INSERT INTO gen(
	user_id,
	title,
	description,
	access,
	body,
	body_meta,
	date_added,
	date_updated
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;