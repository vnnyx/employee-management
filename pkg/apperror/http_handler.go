package apperror

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/invopop/validation"
	"github.com/pkg/errors"
	"github.com/vnnyx/employee-management/internal/constants"
	errorsmap "github.com/vnnyx/employee-management/pkg/errors"
	"github.com/vnnyx/employee-management/pkg/observability/instrumentation"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

type Error struct {
	Issues    []ErrorIssue   `json:"issues,omitempty"`
	Formatted map[string]any `json:"formatted,omitempty"`
}

type ErrorIssue struct {
	Code     string   `json:"issue_code,omitempty"`
	Path     []string `json:"path,omitempty"`
	Message  string   `json:"message,omitempty"`
	Expected any      `json:"expected,omitempty"`
	Received any      `json:"received,omitempty"`
}

func HTTPHandleError(c *fiber.Ctx, err error) error {
	if err == nil {
		return nil
	}

	// AppError
	{
		if appError := (AppError{}); errors.As(err, &appError) {
			if errors.Is(appError.Err, ErrBadRequest) {
				return c.Status(http.StatusBadRequest).JSON(
					Error{
						Issues: []ErrorIssue{
							{
								Code:     appError.IssueCode,
								Path:     appError.Path,
								Message:  appError.Message,
								Expected: appError.Expected,
								Received: appError.Received,
							},
						},
					},
				)
			}

			if errors.Is(appError.Err, ErrForbidden) {
				return c.Status(http.StatusForbidden).JSON(
					Error{
						Issues: []ErrorIssue{
							{
								Code:     appError.IssueCode,
								Path:     appError.Path,
								Message:  appError.Message,
								Expected: appError.Expected,
								Received: appError.Received,
							},
						},
					},
				)
			}

			if errors.Is(appError.Err, ErrNotFound) {
				return c.Status(http.StatusNotFound).JSON(
					Error{
						Issues: []ErrorIssue{
							{
								Code:     appError.IssueCode,
								Path:     appError.Path,
								Message:  appError.Message,
								Expected: appError.Expected,
								Received: appError.Received,
							},
						},
					},
				)
			}
		}
	}

	span := trace.SpanFromContext(c.UserContext())
	var stacks string

	{
		type stackTracer interface {
			StackTrace() errors.StackTrace
		}

		if errStack, ok := errors.Cause(err).(stackTracer); ok {
			st := errStack.StackTrace()
			if len(st) > 4 {
				st = st[0:4]
			}
			stacks = fmt.Sprintf("%+v", st)

			span.AddEvent(
				semconv.ExceptionEventName,
				trace.WithAttributes(
					semconv.ExceptionStacktrace(stacks),
				),
			)
		}
	}

	// JWT errors
	{
		if errors.Is(err, jwt.ErrTokenExpired) {
			return c.Status(http.StatusUnauthorized).JSON(
				Error{
					Issues: []ErrorIssue{
						{
							Code:    constants.AuthTokenExpired,
							Message: errorsmap.GetErrorMessageByIssueCode(constants.AuthTokenExpired),
						},
					},
				},
			)
		}

		if errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return c.Status(http.StatusUnauthorized).JSON(
				Error{
					Issues: []ErrorIssue{
						{
							Code:    constants.AuthTokenInvalid,
							Message: errorsmap.GetErrorMessageByIssueCode(constants.AuthTokenInvalid),
						},
					},
				},
			)
		}

		if errors.Is(err, jwt.ErrTokenMalformed) {
			return c.Status(http.StatusUnauthorized).JSON(
				Error{
					Issues: []ErrorIssue{
						{
							Code:    constants.AuthTokenMalformed,
							Message: errorsmap.GetErrorMessageByIssueCode(constants.AuthTokenMalformed),
						},
					},
				},
			)
		}

		if errors.Is(err, jwt.ErrTokenUnverifiable) {
			return c.Status(http.StatusUnauthorized).JSON(
				Error{
					Issues: []ErrorIssue{
						{
							Code:    constants.AuthTokenNotFound,
							Message: errorsmap.GetErrorMessageByIssueCode(constants.AuthTokenNotFound),
						},
					},
				},
			)
		}
	}

	// Validation errors
	{
		if validationError := (validation.Errors{}); errors.As(err, &validationError) {
			return c.Status(http.StatusBadRequest).JSON(
				validationErrorMapping(validationError),
			)
		}
	}

	// Debug Mode for local env
	var recoveredStack []byte
	if v, ok := c.Locals("recoveredStack").([]byte); ok {
		recoveredStack = v
	}
	isLocal := os.Getenv("APP_ENV") == "local"
	if isLocal {
		log.Println(err)
		if recoveredStack != nil {
			log.Println(string(recoveredStack))
		} else if len(stacks) != 0 {
			log.Println(stacks)
		}
	}

	if recoveredStack != nil {
		instrumentation.RecordSpanErrorWithStackTrace(span, err)
	}

	if isLocal {
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	// Default Fiber Error
	if fiberErr := new(fiber.Error); errors.As(err, &fiberErr) {
		return c.SendStatus(fiberErr.Code)
	}

	return c.Status(http.StatusInternalServerError).JSON(
		Error{
			Issues: []ErrorIssue{
				{
					Code:    constants.InternalServerError,
					Message: errorsmap.GetErrorMessageByIssueCode(constants.InternalServerError),
				},
			},
		},
	)
}

func validationErrorMapping(validationError validation.Errors) Error {
	flatten := validationErrorMappingToFlatten(validationError)

	issues := make([]ErrorIssue, 0, len(flatten))
	for key, errs := range flatten {
		issues = append(issues, ErrorIssue{
			Code:    constants.ValidationError,
			Path:    strings.Split(key, "."),
			Message: strings.Join(errs, ", "),
		})
	}

	return Error{
		Issues:    issues,
		Formatted: fromFlattenToFormatted(flatten),
	}
}

func validationErrorMappingToFlatten(validatorError validation.Errors) map[string][]string {
	mapErr := make(map[string][]string)

	var flatten func(key string, err validation.Errors)
	flatten = func(key string, err validation.Errors) {
		for k, v := range err {
			fullKey := k
			if key != "" {
				fullKey = key + "." + k
			}

			if nestedErrs, ok := v.(validation.Errors); ok {
				flatten(fullKey, nestedErrs)
			} else {
				mapErr[fullKey] = append(mapErr[fullKey], v.Error())
			}
		}
	}

	flatten("", validatorError)
	return mapErr
}

func fromFlattenToFormatted(errFlatten map[string][]string) map[string]any {
	formatted := make(map[string]any)

	for key, messages := range errFlatten {
		path := strings.Split(key, ".")
		current := formatted

		for i, segment := range path {
			if i == len(path)-1 {
				if _, exists := current[segment]; !exists {
					current[segment] = map[string]any{"_errors": messages}
				}
				break
			}

			if _, exists := current[segment]; !exists {
				current[segment] = make(map[string]any)
			}

			current = current[segment].(map[string]any)
		}
	}

	return formatted
}
