package httpController

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Arup3201/gotasks/internal/errors"
)

const (
	NOOP         = "NO_MODIFICATION"
	MISSINGBODY  = "MISSING_BODY_PROPERTY"
	INVALIDBODY  = "INVALID_BODY_PROPERTY"
	INVALIDPARAM = "INVALID_PARAMETER_VALUE"
	NOTFOUND     = "NOT_FOUND"
	SERVERERROR  = "SERVER_ERROR"
)

type BaseError struct {
	Id     string `json:"id"`
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

type HttpError struct {
	BaseError
	Errors []ErrorField `json:"errors"`
}

func (e *HttpError) Error() string {
	if e.Cause == nil {
		return e.Detail
	}

	return e.Detail + ": " + e.Cause.Error()
}

func (e *HttpError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}

	return body, nil
}

func (e *HttpError) ResponseHeader() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/problem+json; charset=utf-8",
	}
}

func New(errId string, errType string, title string, detail string, status int, code string, cause error, fields ...ErrorField) *HttpError {
	return &HttpError{
		BaseError: BaseError{
			Id:     errId,
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

func FromAppError(appError *errors.AppError) *HttpError {
	if appError.Type == errors.INVALID_INPUT {
		fields := []ErrorField{}
		for _, errorField := range appError.Errors {
			fields = append(fields, ErrorField{
				Field:  errorField.Field,
				Reason: errorField.Reason,
			})
		}
		return InvalidBodyError(fields...)
	}
	if appError.Type == errors.NOT_FOUND {
		return NotFoundError()
	}
	if appError.Type == errors.NO_OPERATION {
		return NoOpError()
	}

	return InternalServerError(appError.Cause)
}

func InvalidBodyError(fields ...ErrorField) *HttpError {
	return New(
		INVALIDBODY,
		"https://problems-registry.smartbear.com/invalid-body-property-value",
		"Invalid body property value",
		"The request body contains an invalid body property value.",
		http.StatusBadRequest,
		"400-07",
		nil,
		fields...,
	)
}

func MissingBodyError(fields ...ErrorField) *HttpError {
	return New(
		MISSINGBODY,
		"https://problems-registry.smartbear.com/missing-body-property",
		"Missing body property",
		"The request is missing an expected body property.",
		http.StatusBadRequest,
		"400-09",
		nil,
		fields...,
	)
}

func InvalidRequestParamError(fields ...ErrorField) *HttpError {
	return New(
		INVALIDPARAM,
		"https://problems-registry.smartbear.com/invalid-request-parameter-value",
		"Invalid request parameter value",
		"The request body contains an invalid request parameter value.",
		http.StatusBadRequest,
		"400-08",
		nil,
		fields...,
	)
}

func NotFoundError() *HttpError {
	return New(
		NOTFOUND,
		"https://problems-registry.smartbear.com/not-found",
		"Not found",
		"The requested resource was not found",
		http.StatusNotFound,
		"404-01",
		nil,
	)
}

func NoOpError() *HttpError {
	return New(
		NOOP,
		"about:blank",
		"Not modified",
		"No modification happened at server",
		http.StatusNotModified,
		"204-01",
		nil,
	)
}

func InternalServerError(cause error) *HttpError {
	return New(
		SERVERERROR,
		"https://problems-registry.smartbear.com/server-error",
		"Internal server error",
		"The server encountered an unexpected error",
		http.StatusInternalServerError,
		"500-01",
		cause,
	)
}
