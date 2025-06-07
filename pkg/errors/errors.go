package errors

import "github.com/vnnyx/employee-management/internal/constants"

func GetErrorMessageByIssueCode(issueCode string) string {
	switch issueCode {
	case constants.AuthTokenExpired:
		return "Token has expired, please login again"
	case constants.AuthTokenInvalid:
		return "Token is invalid, please login again"
	case constants.AuthTokenMalformed:
		return "Token is malformed, please login again"
	case constants.AuthTokenNotFound:
		return "Token not found, please login again"
	}

	return "An unknown error occurred"
}
