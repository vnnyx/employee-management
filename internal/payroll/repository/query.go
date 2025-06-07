package repository

const insertPayrollQuery = `
INSERT INTO payrolls (
	id,
	period_id,
	run_by,
	run_at,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
)
VALUES (
	:id,
	:period_id,
	:run_by,
	:run_at,
	:created_at,
	:updated_at,
	:created_by,
	:updated_by,
	:ip_address
)
RETURNING id
`

const insertPayslipQuery = `
INSERT INTO payslips (
	id,
	user_id,
	payroll_id,
	base_salary,
	attendance_days,
	overtime_hours,
	overtime_pay,
	reimbursement_total,
	total_take_home,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
)
VALUES (
	:id,
	:user_id,
	:payroll_id,
	:base_salary,
	:attendance_days,
	:overtime_hours,
	:overtime_pay,
	:reimbursement_total,
	:total_take_home,
	:created_at,
	:updated_at,
	:created_by,
	:updated_by,
	:ip_address
)
RETURNING id
`

const insertPayrollSummaryQuery = `
INSERT INTO payroll_summaries (
	id,
	payroll_id,
	total_take_home,
	generated_by,
	generated_at,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
)
VALUES (
	:id,
	:payroll_id,
	:total_take_home,
	:generated_by,
	:generated_at,
	:created_at,
	:updated_at,
	:created_by,
	:updated_by,
	:ip_address
)
RETURNING id
`

const findPayrollByPeriodIDQuery = `
SELECT
	id,
	period_id,
	run_by,
	run_at,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
FROM payrolls
WHERE period_id = $1
`

const findPayslipByUserIDPeriodQuery = `
SELECT
	id,
	user_id,
	payroll_id,
	base_salary,
	attendance_days,
	overtime_hours,
	overtime_pay,
	reimbursement_total,
	total_take_home,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
FROM payslips
WHERE user_id = $1 AND payroll_id = $2
`

const findPayslipByPayrollIDQuery = `
SELECT
	id,
	user_id,
	payroll_id,
	base_salary,
	attendance_days,
	overtime_hours,
	overtime_pay,
	reimbursement_total,
	total_take_home,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
FROM payslips
WHERE payroll_id = $1
`

const findPayrollByIDQuery = `
SELECT
	id,
	period_id,
	run_by,
	run_at,
	created_at,
	updated_at,
	created_by,
	updated_by,
	ip_address
FROM payrolls
WHERE id = $1
`
