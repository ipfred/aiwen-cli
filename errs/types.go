package errs

import "fmt"

const (
	TypeValidation = "validation"
	TypeConfig     = "config"
	TypeAPIError   = "api_error"
	TypeNetwork    = "network"
	TypeParseError = "parse_error"
	TypeInternal   = "internal"
)

const (
	ExitInternal   = 1
	ExitValidation = 2
	ExitConfig     = 3
	ExitAPI        = 4
	ExitNetwork    = 5
)

type Problem struct {
	Type    string         `json:"type"`
	Message string         `json:"message"`
	Hint    string         `json:"hint,omitempty"`
	Detail  map[string]any `json:"detail,omitempty"`
}

type ExitError struct {
	Code    int
	Problem Problem
}

func (e *ExitError) Error() string {
	return e.Problem.Message
}

func New(code int, typ, message string) *ExitError {
	return &ExitError{
		Code: code,
		Problem: Problem{
			Type:    typ,
			Message: message,
		},
	}
}

func WithDetail(code int, typ, message string, detail map[string]any) *ExitError {
	err := New(code, typ, message)
	err.Problem.Detail = detail
	return err
}

func Validation(message string) *ExitError {
	return New(ExitValidation, TypeValidation, message)
}

func Validationf(format string, args ...any) *ExitError {
	return Validation(fmt.Sprintf(format, args...))
}

func Config(message string) *ExitError {
	return New(ExitConfig, TypeConfig, message)
}

func API(message string) *ExitError {
	return New(ExitAPI, TypeAPIError, message)
}

func Network(message string) *ExitError {
	return New(ExitNetwork, TypeNetwork, message)
}

func Parse(message string) *ExitError {
	return New(ExitAPI, TypeParseError, message)
}
