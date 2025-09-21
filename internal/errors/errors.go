package errors

const (
	INVALID_INPUT = "INVALID_INPUT"
	NOT_FOUND     = "NOT_FOUND"
	NO_OPERATION  = "NOOP"
)

type AppError struct {
	Type   string
	Title  string
	Detail string
	Cause  error
}

func New(type_, title, detail string, cause error) *AppError {
	return &AppError{
		Type:   type_,
		Title:  title,
		Detail: detail,
		Cause:  cause,
	}
}

func (e *AppError) Error() string {
	return e.Detail
}

func InputValidationError(title, detail string) *AppError {
	return New(INVALID_INPUT, title, detail, nil)
}

func NotFoundError(detail string) *AppError {
	return New(NOT_FOUND, "Resource not found", detail, nil)
}

func NoOp(detail string) *AppError {
	return New(NO_OPERATION, "No operation happened", detail, nil)
}
