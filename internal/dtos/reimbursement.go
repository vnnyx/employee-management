package dtos

import (
	"time"

	"github.com/invopop/validation"
	"github.com/vnnyx/employee-management/internal/reimbursement/entity"
	"github.com/vnnyx/employee-management/pkg/optional"
)

type ReimbursementRequest struct {
	Amount      int64           `json:"amount" validate:"required"`
	Date        string          `json:"date" validate:"required"`
	Description optional.String `json:"description,omitempty"`
}

func (r *ReimbursementRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Amount, validation.Required, validation.Min(1)),
		validation.Field(&r.Date, validation.Required, validation.Date("2006-01-02")),
	)
}

func (r *ReimbursementRequest) ToRequestEntity() entity.SubmitReimbursement {
	parsedDate, _ := time.Parse("2006-01-02", r.Date)
	return entity.SubmitReimbursement{
		Amount:      r.Amount,
		Date:        parsedDate,
		Description: r.Description,
	}
}
