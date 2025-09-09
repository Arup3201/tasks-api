package apperr

const (
	BadRequest = 400
	NotFound   = 404
)

func InvalidBodyError(fields ...ErrorField) *Error {
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

func MissingBodyError(fields ...ErrorField) *Error {
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

func InvalidRequestParamError(fields ...ErrorField) *Error {
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

func NotFoundError() *Error {
	return New(
		"https://problems-registry.smartbear.com/not-found",
		"Not Found",
		"The requested resource was not found",
		NotFound,
		"404-01",
		nil,
	)
}
