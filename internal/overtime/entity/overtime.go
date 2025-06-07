package entity

import (
	"time"
)

type Overtime struct {
	ID            string        `db:"id"`
	UserID        string        `db:"user_id"`
	OverTimeDate  time.Time     `db:"overtime_date"`
	OvertimeHours time.Duration `db:"overtime_hours"`
	CreatedAt     time.Time     `db:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at"`
	CreatedBy     string        `db:"created_by"`
	UpdatedBy     string        `db:"updated_by"`
	IPAddress     string        `db:"ip_address"`
}

type SubmitOvertime struct {
	OvertimeDate time.Time
	Overtime     string
}

type MappedBy string

const (
	MappedByUserID         MappedBy = "user_id"
	MappedByAttendanceDate MappedBy = "overtime_date"
)

type MappedOptions struct {
	MappedBy MappedBy
}

type FindOvertimeOptions struct {
	PessimisticLock bool
	*MappedOptions
}

type FindOvertimeResult struct {
	List     []Overtime
	Mapped   map[any][]Overtime
	IsMapped bool
	MappedBy MappedBy
}
