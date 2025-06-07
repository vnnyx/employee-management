package entity

const (
	PayrollNotAuthorized     = "PAYROLL_NOT_AUTHORIZED"
	PayrollAlreadyGenerated  = "PAYROLL_ALREADY_GENERATED"
	AttendancePeriodNotFound = "ATTENDANCE_PERIOD_NOT_FOUND"
	UserNotFound             = "USER_NOT_FOUND"
	PayslipNotFound          = "PAYSLIP_NOT_FOUND"
)

func GetErrorMessageByIssueCode(issueCode string) string {
	switch issueCode {
	case PayrollNotAuthorized:
		return "You are not authorized to perform this action"
	case PayrollAlreadyGenerated:
		return "Payroll for this period has already been generated"
	case AttendancePeriodNotFound:
		return "Attendance period not found"
	case UserNotFound:
		return "User not found"
	case PayslipNotFound:
		return "Payslip not found"
	default:
		return "An unknown error occurred"
	}
}
