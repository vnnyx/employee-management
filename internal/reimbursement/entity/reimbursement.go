package entity

import (
	"time"

	"github.com/vnnyx/employee-management/pkg/optional"
)

type Reimbursement struct {
	ID                string          `db:"id"`
	UserID            string          `db:"user_id"`
	Amount            int64           `db:"amount"`
	Description       optional.String `db:"description"`
	ReimbursementDate time.Time       `db:"reimbursement_date"`
	CreatedAt         time.Time       `db:"created_at"`
	UpdatedAt         time.Time       `db:"updated_at"`
	CreatedBy         string          `db:"created_by"`
	UpdatedBy         string          `db:"updated_by"`
	IPAddress         string          `db:"ip_address"`
}

type SubmitReimbursement struct {
	Amount      int64
	Date        time.Time
	Description optional.String
}

type MappedBy string

const (
	MappedByUserID MappedBy = "user_id"
	MappedByDate   MappedBy = "reimbursement_date"
)

type MappedOptions struct {
	MappedBy MappedBy
}

type FindReimbursementOptions struct {
	PessimisticLock bool
	*MappedOptions
}

type FindReimbursementResult struct {
	List     []Reimbursement
	Mapped   map[any][]Reimbursement
	IsMapped bool
	MappedBy MappedBy
}
