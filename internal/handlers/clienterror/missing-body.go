package clienterror

import (
	"encoding/json"
	"fmt"
)

type MissingBodyPropertyError struct {
	BaseError
	Errors []ErrorDetail `json:"errors"`
}

func (e *MissingBodyPropertyError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}

	return e.Detail + ": " + e.Cause.Error()
}

func (e *MissingBodyPropertyError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}

	return body, nil
}

func (e *MissingBodyPropertyError) ResponseHeader() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/problem+json; charset=utf-8",
	}
}

func NewMissingBodyProperyError(err error, errors []ErrorDetail) error {
	return &MissingBodyPropertyError{
		BaseError: BaseError{
			Type:   "https://problems-registry.smartbear.com/missing-body-property",
			Title:  "Missing body property",
			Detail: "The request is missing an expected body property.",
			Status: 400,
			Code:   "400-09",
			Cause:  err,
		},
		Errors: errors,
	}
}
