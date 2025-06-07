package repository

const insertAttendanceQuery = `
INSERT INTO attendances (
	id,
	user_id,
	attendance_date,
	created_at,
	updated_at,
	updated_by,
	created_by,
	ip_address
)
VALUES (
	:id,
	:user_id,
	:attendance_date,
	:created_at,
	:updated_at,
	:updated_by,
	:created_by,
	:ip_address
)
RETURNING id
`

const upsertAttendanceQuery = `
INSERT INTO attendances (
	id,
	user_id,
	attendance_date,
	created_at,
	updated_at,
	updated_by,
	created_by,
	ip_address
)
VALUES (
	:id,
	:user_id,
	:attendance_date,
	:created_at,
	:updated_at,
	:updated_by,
	:created_by,
	:ip_address
)
ON CONFLICT (user_id, attendance_date) DO UPDATE SET
	updated_at = EXCLUDED.updated_at,
	updated_by = EXCLUDED.updated_by,
	ip_address = EXCLUDED.ip_address
RETURNING id
`

const insertAttendancePeriodQuery = `
INSERT INTO attendance_periods (
	id,
	start_date,
	end_date,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
)
VALUES (
	:id,
	:start_date,
	:end_date,
	:created_at,
	:updated_at,
	:created_by,
	:updated_by,
	:ip_address
)
RETURNING id
`

const findAttendancePeriodByIDQuery = `
SELECT
	id,
	start_date,
	end_date,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
FROM attendance_periods
WHERE id = $1
`

const findAttendanceByPeriodQuery = `
SELECT
	id,
	user_id,
	attendance_date,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
FROM attendances
WHERE attendance_date BETWEEN $1::DATE AND $2::DATE
ORDER BY attendance_date
`

const findAttendancePeriodByPayrollIDQuery = `
SELECT
	ad.id,
	ad.start_date,
	ad.end_date,
	ad.created_at,
	ad.updated_at,
	ad.created_by,
	ad.updated_by,
	ad.ip_address
FROM attendance_periods ad
JOIN payrolls p ON ad.id = p.period_id
WHERE p.id = $1
`
