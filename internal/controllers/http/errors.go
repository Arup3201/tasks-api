package httpController

import (
	"encoding/json"
	"fmt"

	"github.com/Arup3201/gotasks/internal/errors"
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

func New(errType string, title string, detail string, status int, code string, cause error, fields ...ErrorField) *HttpError {
	return &HttpError{
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

func FromAppError(appError *errors.AppError) *HttpError {
	if appError.Type == errors.INVALID_INPUT {
		return InvalidBodyError()
	}
	if appError.Type == errors.NOT_FOUND {
		return NotFoundError()
	}
	if appError.Type == errors.NO_OPERATION {
		return NoOpError()
	}

	return InternalServerError(appError.Cause)
}

// helpers
const (
	NoOp        = 204
	BadRequest  = 400
	NotFound    = 404
	ServerError = 500
)

func InvalidBodyError(fields ...ErrorField) *HttpError {
	return New(
		"https://problems-registry.smartbear.com/invalid-body-property-value",
		"Invalid Body Property Value",
		"The request body contains an invalid body property value.",
		BadRequest,
		"400-07",
		nil,
		fields...,
	)
}

func MissingBodyError(fields ...ErrorField) *HttpError {
	return New(
		"https://problems-registry.smartbear.com/missing-body-property",
		"Missing body property",
		"The request is missing an expected body property.",
		BadRequest,
		"400-09",
		nil,
		fields...,
	)
}

func InvalidRequestParamError(fields ...ErrorField) *HttpError {
	return New(
		"https://problems-registry.smartbear.com/invalid-request-parameter-value",
		"Invalid Request Parameter Value",
		"The request body contains an invalid request parameter value.",
		BadRequest,
		"400-08",
		nil,
		fields...,
	)
}

func NotFoundError() *HttpError {
	return New(
		"https://problems-registry.smartbear.com/not-found",
		"Not Found",
		"The requested resource was not found",
		NotFound,
		"404-01",
		nil,
	)
}

func NoOpError() *HttpError {
	return New(
		"about:blank",
		"Not Modified",
		"No modification happened at server",
		NoOp,
		"204-01",
		nil,
	)
}

func InternalServerError(cause error) *HttpError {
	return New(
		"https://problems-registry.smartbear.com/server-error",
		"Server Error",
		"The server encountered an unexpected error",
		ServerError,
		"500-01",
		cause,
	)
}
