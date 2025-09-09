package apperr

type ClientError interface {
	Error() string
	// ResponseBody returns response body of the error
	ResponseBody() ([]byte, error)
	// ResponseHeader returns status code and response headers of the error
	ResponseHeader() (int, map[string]string)
}
