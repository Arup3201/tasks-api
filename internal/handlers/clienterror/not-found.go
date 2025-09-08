package clienterror

import (
	"encoding/json"
	"fmt"
)

type NotFoundError struct {
	BaseError
}

func (e *NotFoundError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}

	return e.Detail + ": " + e.Cause.Error()
}

func (e *NotFoundError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}

	return body, nil
}

func (e *NotFoundError) ResponseHeader() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/problem+json; charset=utf-8",
	}
}

func NewNotFoundError(err error) error {
	return &NotFoundError{
		BaseError: BaseError{
			Type:   "https://problems-registry.smartbear.com/not-found",
			Title:  "Not Found",
			Detail: "The requested resource was not found",
			Status: 404,
			Code:   "404-01",
			Cause:  err,
		},
	}
}
