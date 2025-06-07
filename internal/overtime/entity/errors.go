package entity

const (
	OvertimeInvalidTimeRequest = "OVERTIME_INVALID_TIME_REQUEST"
	OvertimeExceedsLimit       = "OVERTIME_EXCEEDS_LIMIT"
)

func GetErrorMessageByIssueCode(issueCode string) string {
	switch issueCode {
	case OvertimeInvalidTimeRequest:
		return "Overtime cannot be submitted on working hours"
	case OvertimeExceedsLimit:
		return "Overtime exceeds the allowed limit for the day"
	default:
		return "An unknown error occurred"
	}
}
