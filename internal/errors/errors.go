package errors

type ClientError interface {
	Error() string
	// ResponseBody returns response body of the error
	ResponseBody() ([]byte, error)
	// ResponseHeader returns status code and response headers of the error
	ResponseHeader() (int, map[string]string)
}

type BaseError struct {
	Type   string `json:"type"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
	Code   string `json:"code"`
	Cause  error  `json:"-"`
}
