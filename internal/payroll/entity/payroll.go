package entity

import (
	"time"

	"github.com/vnnyx/employee-management/pkg/optional"
	// reimbursementEntity "github.com/vnnyx/employee-management/internal/reimbursement/entity"
	// overtimeEntity "github.com/vnnyx/employee-management/internal/overtime/entity"
)

type Payroll struct {
	ID        string    `db:"id"`
	PeriodID  string    `db:"period_id"`
	RunBy     string    `db:"run_by"`
	RunAt     time.Time `db:"run_at"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedBy string    `db:"created_by"`
	UpdatedBy string    `db:"updated_by"`
	IPAddress string    `db:"ip_address"`
}

type Payslip struct {
	ID                 string            `db:"id"`
	UserID             string            `db:"user_id"`
	PayrollID          string            `db:"payroll_id"`
	BaseSalary         int64             `db:"base_salary"`
	AttendanceDays     int64             `db:"attendance_days"`
	OvertimeHours      optional.Duration `db:"overtime_hours"`
	OvertimePay        int64             `db:"overtime_pay"`
	ReimbursementTotal int64             `db:"reimbursement_total"`
	TotalTakeHome      int64             `db:"total_take_home"`
	CreatedAt          time.Time         `db:"created_at"`
	UpdatedAt          time.Time         `db:"updated_at"`
	CreatedBy          string            `db:"created_by"`
	UpdatedBy          string            `db:"updated_by"`
	IPAddress          string            `db:"ip_address"`
}

type PayrollSummary struct {
	ID            string    `db:"id"`
	PayrollID     string    `db:"payroll_id"`
	TotalTakeHome int64     `db:"total_take_home"`
	GeneratedBy   string    `db:"generated_by"`
	GeneratedAt   time.Time `db:"generated_at"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	CreatedBy     string    `db:"created_by"`
	UpdatedBy     string    `db:"updated_by"`
	IPAddress     string    `db:"ip_address"`
}

type GeneratedPayroll struct {
	PeriodID      string
	PayrollID     string
	TotalTakeHome int64
	TotalEmployee int64
	TotalPayslip  int64
	GeneratedBy   string
	GeneratedAt   string
}

type FindPayrollOptions struct {
	PessimisticLock bool
}

type ReimbursementData struct {
	Description       optional.String
	Amount            int64
	ReimbursementDate string
}

type OvertimeData struct {
	OvertimeHours string
	RatePerHour   int64
	Multiplier    float64
	OvertimePay   int64
}

type UserData struct {
	ID       string
	Username string
}

type AttendancePeriodData struct {
	StartDate string
	EndDate   string
}

type PayslipData struct {
	ID                 string
	User               UserData
	AttendancePeriod   AttendancePeriodData
	BaseSalary         int64
	WorkingDays        int64
	AttendanceDays     int64
	AttendancePay      int64
	Overtime           OvertimeData
	Reimbursements     []ReimbursementData
	ReimbursementTotal int64
	TotalTakeHome      int64
}

type ListPayslips struct {
	TotalTakeHome int64
	PayslipsData  []PayslipData
}

type MappedBy string

const (
	MappedByUserID    MappedBy = "user_id"
	MappedByPayslipID MappedBy = "payslip_id"
	MappedByPayrollID MappedBy = "payroll_id"
	MapByPeriodID     MappedBy = "period_id"
)

type MappedOptions struct {
	MappedBy MappedBy
}

type FindPayslipOptions struct {
	PessimisticLock bool
	*MappedOptions
}

type FindPayslipResult struct {
	List     []Payslip
	Mapped   map[any][]Payslip
	IsMapped bool
	MappedBy MappedBy
}
