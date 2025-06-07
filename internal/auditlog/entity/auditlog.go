package entity

import "encoding/json"

type AuditLog struct {
	ID        string          `db:"id"`
	TableName string          `db:"table_name"`
	RecordID  string          `db:"record_id"`
	Action    string          `db:"action"`
	ChangedBy string          `db:"changed_by"`
	IPAddress string          `db:"ip_address"`
	RequestID string          `db:"request_id"`
	OldData   json.RawMessage `db:"old_data"`
	NewData   json.RawMessage `db:"new_data"`
	CreatedAt string          `db:"created_at"`
}
