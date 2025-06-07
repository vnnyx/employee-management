package dtos

import (
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	"github.com/vnnyx/employee-management/internal/payroll/entity"
	"github.com/vnnyx/employee-management/pkg/optional"
)

type GeneratedPayrollResponse struct {
	PeriodID      string `json:"period_id"`
	PayrollID     string `json:"payroll_id"`
	TotalTakeHome int64  `json:"total_take_home_pay"`
	TotalEmployee int64  `json:"total_employee"`
	TotalPayslip  int64  `json:"total_payslip"`
	GeneratedBy   string `json:"generated_by"`
	GeneratedAt   string `json:"generated_at"`
}

type GeneratePayrollRequest struct {
	PeriodID string `json:"period_id" validate:"required"`
}

func (r *GeneratePayrollRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.PeriodID, validation.Required, is.UUIDv4),
	)
}

type ReimbursementDataResponse struct {
	Description       optional.String `json:"description"`
	Amount            int64           `json:"amount"`
	ReimbursementDate string          `json:"reimbursement_date"`
}

type OvertimeDataResponse struct {
	OvertimeHours string  `json:"overtime_hours"`
	RatePerHour   int64   `json:"rate_per_hour"`
	Multiplier    float64 `json:"multiplier"`
	OvertimePay   int64   `json:"overtime_pay"`
}

type UserDataResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type AttendancePeriodDataResponse struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type PayslipDataResponse struct {
	ID                 string                       `json:"id"`
	User               UserDataResponse             `json:"user"`
	AttendancePeriod   AttendancePeriodDataResponse `json:"attendance_period"`
	BaseSalary         int64                        `json:"base_salary"`
	WorkingDays        int64                        `json:"working_days"`
	AttendanceDays     int64                        `json:"attendance_days"`
	AttendancePay      int64                        `json:"attendance_pay"`
	Overtime           OvertimeDataResponse         `json:"overtime"`
	Reimbursements     []ReimbursementDataResponse  `json:"reimbursements"`
	ReimbursementTotal int64                        `json:"reimbursement_total"`
	TotalTakeHome      int64                        `json:"total_take_home_pay"`
}

func NewShowPayslipResponse(data *entity.PayslipData) *PayslipDataResponse {
	return &PayslipDataResponse{
		ID:                 data.ID,
		User:               UserDataResponse{ID: data.User.ID, Username: data.User.Username},
		AttendancePeriod:   AttendancePeriodDataResponse{StartDate: data.AttendancePeriod.StartDate, EndDate: data.AttendancePeriod.EndDate},
		BaseSalary:         data.BaseSalary,
		WorkingDays:        data.WorkingDays,
		AttendanceDays:     data.AttendanceDays,
		AttendancePay:      data.AttendancePay,
		Overtime:           OvertimeDataResponse{OvertimeHours: data.Overtime.OvertimeHours, RatePerHour: data.Overtime.RatePerHour, Multiplier: data.Overtime.Multiplier, OvertimePay: data.Overtime.OvertimePay},
		Reimbursements:     newReimbursementDataResponses(data.Reimbursements),
		ReimbursementTotal: data.ReimbursementTotal,
		TotalTakeHome:      data.TotalTakeHome,
	}
}

type ListPayslipResponse struct {
	TotalTakeHome int64
	PayslipsData  []PayslipDataResponse
}

func NewListPayslipsResponse(data *entity.ListPayslips) *ListPayslipResponse {
	payslips := make([]PayslipDataResponse, len(data.PayslipsData))
	for i, p := range data.PayslipsData {
		payslips[i] = PayslipDataResponse{
			ID:                 p.ID,
			User:               UserDataResponse{ID: p.User.ID, Username: p.User.Username},
			AttendancePeriod:   AttendancePeriodDataResponse{StartDate: p.AttendancePeriod.StartDate, EndDate: p.AttendancePeriod.EndDate},
			BaseSalary:         p.BaseSalary,
			WorkingDays:        p.WorkingDays,
			AttendanceDays:     p.AttendanceDays,
			AttendancePay:      p.AttendancePay,
			Overtime:           OvertimeDataResponse{OvertimeHours: p.Overtime.OvertimeHours, RatePerHour: p.Overtime.RatePerHour, Multiplier: p.Overtime.Multiplier, OvertimePay: p.Overtime.OvertimePay},
			Reimbursements:     newReimbursementDataResponses(p.Reimbursements),
			ReimbursementTotal: p.ReimbursementTotal,
			TotalTakeHome:      p.TotalTakeHome,
		}
	}

	return &ListPayslipResponse{
		TotalTakeHome: data.TotalTakeHome,
		PayslipsData:  payslips,
	}
}

func newReimbursementDataResponses(reimbursements []entity.ReimbursementData) []ReimbursementDataResponse {
	reimbursementResponses := make([]ReimbursementDataResponse, len(reimbursements))
	for i, r := range reimbursements {
		reimbursementResponses[i] = ReimbursementDataResponse{
			Description:       r.Description,
			Amount:            r.Amount,
			ReimbursementDate: r.ReimbursementDate,
		}
	}
	return reimbursementResponses
}
