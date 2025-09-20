package errors

type Error struct {
	Type   string
	Title  string
	Detail string
}

func New(type_, title, detail string) Error {
	return Error{
		Type:   type_,
		Title:  title,
		Detail: detail,
	}
}

func (e Error) Error() string {
	return e.Detail
}

func InputValidationError(title, detail string) Error {
	return New("INVALID_INPUT", title, detail)
}

func NotFoundError(detail string) Error {
	return New("NOT_FOUND", "Resource not found", detail)
}
