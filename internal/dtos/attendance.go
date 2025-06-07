package dtos

import (
	"time"

	"github.com/invopop/validation"
	"github.com/vnnyx/employee-management/internal/attendance/entity"
)

const dateFormat = "2006-01-02"

type AttendancePeriodRequest struct {
	StartDate string `json:"start_date" validate:"required"`
	EndDate   string `json:"end_date" validate:"required"`
}

func (r *AttendancePeriodRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.StartDate, validation.Required, validation.Date(dateFormat)),
		validation.Field(&r.EndDate, validation.Required, validation.Date(dateFormat)),
	)
}

func (r *AttendancePeriodRequest) ToRequestEntity() entity.CreateAttendancePeriod {
	parseStartDate, err := time.Parse(dateFormat, r.StartDate)
	if err != nil {
		panic("invalid start date format, expected 'YYYY-MM-DD'")
	}
	parseEndDate, err := time.Parse(dateFormat, r.EndDate)
	if err != nil {
		panic("invalid end date format, expected 'YYYY-MM-DD'")
	}

	return entity.CreateAttendancePeriod{
		StartDate: parseStartDate,
		EndDate:   parseEndDate,
	}
}
