package repository

const insertReimbursementQuery = `
INSERT INTO reimbursements (
	id,
	user_id,
	amount,
	description,
	reimbursement_date,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
)
VALUES (
	:id,
	:user_id,
	:amount,
	:description,
	:reimbursement_date,
	:created_at,
	:updated_at,
	:created_by,
	:updated_by,
	:ip_address
)
RETURNING id
`

const findReimbursementByPeriodQuery = `
SELECT
	id,
	user_id,
	amount,
	description,
	reimbursement_date,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
FROM reimbursements
WHERE reimbursement_date BETWEEN $1::DATE AND $2::DATE
`

const findReimbursementByUserIDPeriodQuery = `
SELECT
	id,
	user_id,
	amount,
	description,
	reimbursement_date,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
FROM reimbursements
WHERE user_id = $1 AND reimbursement_date BETWEEN $2::DATE AND $3::DATE
`
