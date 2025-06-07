package apperror

import "errors"

var (
	ErrBadRequest   = errors.New("bad_request")
	ErrNotFound     = errors.New("not_found")
	ErrInternal     = errors.New("internal_error")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type AppError struct {
	Err       error
	IssueCode string
	Path      []string
	Message   string
	Expected  any
	Received  any
}

type AppErrors []AppError

func (e AppError) Error() string {
	return e.Err.Error()
}

func (e AppError) Is(err error) bool {
	return errors.Is(e.Err, err)
}

func (e AppError) As(v any) bool {
	return errors.As(e.Err, &v)
}

func BadRequest(appErr AppError) AppError {
	return AppError{
		Err:       ErrBadRequest,
		IssueCode: appErr.IssueCode,
		Path:      appErr.Path,
		Message:   appErr.Message,
		Expected:  appErr.Expected,
		Received:  appErr.Received,
	}
}

func Forbidden(appErr AppError) AppError {
	return AppError{
		Err:       ErrForbidden,
		IssueCode: appErr.IssueCode,
		Path:      appErr.Path,
		Message:   appErr.Message,
		Expected:  appErr.Expected,
		Received:  appErr.Received,
	}
}

func NotFound(appErr AppError) AppError {
	return AppError{
		Err:       ErrNotFound,
		IssueCode: appErr.IssueCode,
		Path: 	appErr.Path,
		Message:  appErr.Message,
		Expected: appErr.Expected,
		Received: appErr.Received,
	}
}
