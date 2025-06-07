package entity

const (
	AttendanceInvalidDay          = "ATTENDANCE_INVALID_DAY"
	AttendanceNotAuthorized       = "ATTENDANCE_NOT_AUTHORIZED"
	AttendanceInvalidPeriod       = "ATTENDANCE_INVALID_PERIOD"
	AttendancePeriodAlreadyExists = "ATTENDANCE_PERIOD_ALREADY_EXISTS"
)

func GetErrorMessageByIssueCode(issueCode string) string {
	switch issueCode {
	case AttendanceInvalidDay:
		return "Attendance cannot be submitted on weekends"
	case AttendanceNotAuthorized:
		return "You are not authorized to perform this action"
	case AttendanceInvalidPeriod:
		return "The attendance period is invalid, start date must be before end date"
	case AttendancePeriodAlreadyExists:
		return "An attendance period with the same start and end date already exists"
	default:
		return "An unknown error occurred"
	}
}
