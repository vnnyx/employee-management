package entity

import (
	"time"
)

type Attendance struct {
	ID             string    `db:"id"`
	UserID         string    `db:"user_id"`
	AttendanceDate time.Time `db:"attendance_date"`
	CreatedAt      time.Time `db:"created_at"`
	UpdatedAt      time.Time `db:"updated_at"`
	CreatedBy      string    `db:"created_by"`
	UpdatedBy      string    `db:"updated_by"`
	IPAddress      string    `db:"ip_address"`
}

type AttendancePeriod struct {
	ID        string    `db:"id"`
	StartDate time.Time `db:"start_date"`
	EndDate   time.Time `db:"end_date"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedBy string    `db:"updated_by"`
	IPAddress string    `db:"ip_address"`
}

type CreateAttendancePeriod struct {
	StartDate time.Time
	EndDate   time.Time
}

type MappedBy string

const (
	MappedByUserID         MappedBy = "user_id"
	MappedByAttendanceDate MappedBy = "attendance_date"
)

type MappedOptions struct {
	MappedBy MappedBy
}

type FindAttendanceOptions struct {
	PessimisticLock bool
	*MappedOptions
}

type FindAttendanceResult struct {
	List     []Attendance
	Mapped   map[any][]Attendance
	IsMapped bool
	MappedBy MappedBy
}
