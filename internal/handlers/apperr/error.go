package apperr

import (
	"encoding/json"
	"fmt"
)

type BaseError struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
	Code   string `json:"code"`
	Cause  error  `json:"-"`
}

type ErrorField struct {
	Reason string `json:"detail"`
	Field  string `json:"field"`
}

type Error struct {
	BaseError
	Errors []ErrorField `json:"errors"`
}

/* ClientError interface implementation */

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Detail
	}

	return e.Detail + ": " + e.Cause.Error()
}

func (e *Error) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}

	return body, nil
}

func (e *Error) ResponseHeader() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/problem+json; charset=utf-8",
	}
}

/* ................................................. */

func New(errType string, title string, detail string, status int, code string, cause error, fields ...ErrorField) *Error {
	return &Error{
		BaseError: BaseError{
			Type:   errType,
			Title:  title,
			Detail: detail,
			Status: status,
			Code:   code,
			Cause:  cause,
		},
		Errors: fields,
	}
}
