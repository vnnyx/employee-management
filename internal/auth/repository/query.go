package repository

const findUserByUsername = `
SELECT
	u.id,
	u.username,
	u.is_admin,
	u.password
FROM users u
WHERE u.username = $1
`
