package repository

const findUsersQuery = `
SELECT
  id,
  username,
  is_admin,
  salary
FROM users
`

const findUserByIDQuery = `
SELECT
  id,
  username,
  is_admin,
  salary
FROM users
WHERE id = $1
`
