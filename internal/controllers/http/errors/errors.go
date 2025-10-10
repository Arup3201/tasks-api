package httperrors

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Arup3201/gotasks/internal/errors"
)

const (
	INCORRECT_CREDENTIAL = "INCORRECT_CREDENTIALS"
	UNAUTHORIZED         = "NOT_AUTHORIZED"
	NO_OP                = "NO_MODIFICATION"
	MISSING_BODY         = "MISSING_BODY_PROPERTY"
	INVALID_BODY         = "INVALID_BODY_PROPERTY"
	INVALID_PARAM        = "INVALID_PARAMETER_VALUE"
	NOT_FOUND            = "NOT_FOUND"
	SERVER_ERROR         = "SERVER_ERROR"
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

func IncorrectCredentialError() *HttpError {
	return New(
		INCORRECT_CREDENTIAL,
		"about:blank",
		"Incorrect credentials for login",
		"Username or password is not valid, please try with correct username and password",
		http.StatusBadRequest,
		"400-07",
		nil,
	)
}

func UnauthorizedError() *HttpError {
	return New(
		UNAUTHORIZED,
		"https://problems-registry.smartbear.com/unauthorized",
		"Unauthorized",
		"Access token not set or invalid, and the requested resource could not be returned",
		http.StatusUnauthorized,
		"401-01",
		nil,
	)
}

func InvalidBodyError(fields ...ErrorField) *HttpError {
	return New(
		INVALID_BODY,
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
		MISSING_BODY,
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
		INVALID_PARAM,
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
		NOT_FOUND,
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
		NO_OP,
		"about:blank",
		"Not modified",
		"No modification happened at server",
		http.StatusNoContent,
		"204-01",
		nil,
	)
}

func InternalServerError(cause error) *HttpError {
	return New(
		SERVER_ERROR,
		"https://problems-registry.smartbear.com/server-error",
		"Internal server error",
		"The server encountered an unexpected error",
		http.StatusInternalServerError,
		"500-01",
		cause,
	)
}
