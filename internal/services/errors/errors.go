package errors

type InputValidationError struct {
	Type   string
	Title  string
	Detail string
}

func NewInputValidationError(title, detail string) InputValidationError {
	return InputValidationError{
		Type:   "INVALID_INPUT",
		Title:  title,
		Detail: detail,
	}
}

func (e InputValidationError) Error() string {
	return e.Detail
}
