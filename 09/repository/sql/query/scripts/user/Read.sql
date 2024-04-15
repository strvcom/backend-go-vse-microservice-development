SELECT
	u.id,
	u.created_at,
	u.updated_at
FROM
	users as u
WHERE
	id = @id
