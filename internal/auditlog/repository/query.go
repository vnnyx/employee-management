package repository

const insertAuditLogQuery = `
INSERT INTO audit_logs (
	id,
	table_name,
	record_id,
	action,
	changed_by,
	ip_address,
	request_id,
	old_data,
	new_data,
	created_at
VALUES (
	:id,
	:table_name,
	:record_id,
	:action,
	:changed_by,
	:ip_address,
	:request_id,
	:old_data,
	:new_data,
	:created_at
)
RETURNING id
`
