SELECT
	u.id,
	u.created_at,
	u.updated_at
FROM
	users as u
ORDER BY u.created_at
