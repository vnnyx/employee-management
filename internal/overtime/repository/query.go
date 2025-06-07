package repository

const insertOvertimeQuery = `
INSERT INTO overtimes (
	id,
	user_id,
	overtime_date,
	overtime_hours,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
)
VALUES (
	:id,
	:user_id,
	:overtime_date,
	:overtime_hours,
	:created_at,
	:updated_at,
	:created_by,
	:updated_by,
	:ip_address
)
RETURNING id
`

const findOvertimeByUserIDDate = `
SELECT 
	id,
	user_id,
	overtime_date,
	overtime_hours,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
FROM overtimes
WHERE user_id = $1 AND overtime_date = $2::DATE
`

const upsertOvertimeQuery = `
INSERT INTO overtimes (
	id,
	user_id,
	overtime_date,
	overtime_hours,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
)
VALUES (
	:id,
	:user_id,
	:overtime_date,
	:overtime_hours,
	:created_at,
	:updated_at,
	:created_by,
	:updated_by,
	:ip_address
)
ON CONFLICT (user_id, overtime_date) DO UPDATE SET
	overtime_hours = EXCLUDED.overtime_hours,
	updated_at = EXCLUDED.updated_at,
	updated_by = EXCLUDED.updated_by,
	ip_address = EXCLUDED.ip_address
RETURNING id
`

const findOvertimeByPeriodQuery = `
SELECT
	id,
	user_id,
	overtime_date,
	overtime_hours,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
FROM overtimes
WHERE overtime_date BETWEEN $1::DATE AND $2::DATE
ORDER BY overtime_date ASC
`