package dtos

import (
	"regexp"
	"time"

	"github.com/invopop/validation"
	"github.com/vnnyx/employee-management/internal/overtime/entity"
)

type OvertimeRequest struct {
	Date     string `json:"date" validate:"required"`
	Overtime string `json:"overtime" validate:"required"`
}

func (o *OvertimeRequest) Validate() error {
	return validation.ValidateStruct(o,
		validation.Field(&o.Date, validation.Required, validation.Date("2006-01-02")),
		validation.Field(&o.Overtime, validation.Required, validation.Match(regexp.MustCompile(`^PT(\d+H)?(\d+M)?$`))),
	)
}

func (o *OvertimeRequest) ToRequestEntity() entity.SubmitOvertime {
	parsedDate, _ := time.Parse("2006-01-02", o.Date)
	return entity.SubmitOvertime{
		OvertimeDate: parsedDate,
		Overtime:     o.Overtime,
	}
}
