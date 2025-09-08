package clienterror

import (
	"encoding/json"
	"fmt"
)

type InvalidBodyValueError struct {
	BaseError
	Errors []ErrorDetail `json:"errors"`
}

func (e *InvalidBodyValueError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}

	return e.Detail + ": " + e.Cause.Error()
}

func (e *InvalidBodyValueError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}

	return body, nil
}

func (e *InvalidBodyValueError) ResponseHeader() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/problem+json; charset=utf-8",
	}
}

func NewInvalidBodyValueError(err error, errors []ErrorDetail) error {
	return &InvalidBodyValueError{
		BaseError: BaseError{
			Type:   "https://problems-registry.smartbear.com/invalid-body-property-value",
			Title:  "Invalid Body Property Value",
			Detail: "The request body contains an invalid body property value.",
			Status: 400,
			Code:   "400-07",
			Cause:  err,
		},
		Errors: errors,
	}
}
