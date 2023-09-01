package errors

type HTTPError struct {
	HTTPStatus int    `example:"400" json:"httpStatus"`
	ErrorKey   string `example:"Bad Request unable to fetch data" json:"errorKey"`
}

// implementing error interface
func (httpError *HTTPError) Error() string {
	return httpError.ErrorKey
}

// newHTTPError creates a new instance of HTTPError
func NewHTTPError(key string, statuscode int) *HTTPError {
	return &HTTPError{
		HTTPStatus: statuscode,
		ErrorKey:   key,
	}
}
