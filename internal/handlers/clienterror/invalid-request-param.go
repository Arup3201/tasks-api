package clienterror

import (
	"encoding/json"
	"fmt"
)

type InvalidRequestParamError struct {
	BaseError
	Errors []RequestParamError `json:"errors"`
}

func (e *InvalidRequestParamError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}

	return e.Detail + ": " + e.Cause.Error()
}

func (e *InvalidRequestParamError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}

	return body, nil
}

func (e *InvalidRequestParamError) ResponseHeader() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/problem+json; charset=utf-8",
	}
}

func NewInvalidRequestParamError(err error, errors []RequestParamError) error {
	return &InvalidRequestParamError{
		BaseError: BaseError{
			Type:   "https://problems-registry.smartbear.com/invalid-request-parameter-value",
			Title:  "Invalid Request Parameter Value",
			Detail: "The request body contains an invalid request parameter value.",
			Status: 400,
			Code:   "400-08",
			Cause:  err,
		},
		Errors: errors,
	}
}
