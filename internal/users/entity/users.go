package entity

import (
	"time"

	"github.com/vnnyx/employee-management/pkg/resourceful"
)

type User struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	IsAdmin   bool      `db:"is_admin"`
	Salary    int64     `db:"salary"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedBy string    `db:"updated_by"`
	IPAddress string    `db:"ip_address"`
}

type MappedBy string

const (
	MappedByUserID MappedBy = "user_id"
)

type MappedOptions struct {
	MappedBy MappedBy
}

type FindUserOptions struct {
	PessimisticLock bool
	*resourceful.CursorParameter
	*MappedOptions
}

type FindUserResult struct {
	List     []User
	Mapped   map[any][]User
	IsMapped bool
	MappedBy MappedBy
}
