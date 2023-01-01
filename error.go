package blia

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrDecodeBodyFailed   = NewHTTPErrorMsg("decode body failed", http.StatusBadRequest)
	ErrDecodeQueryFailed  = NewHTTPErrorMsg("decode query failed", http.StatusBadRequest)
	ErrInvalidOffsetLimit = NewHTTPErrorMsg("invalid offset,limit or page,page_size", http.StatusBadRequest)
	ErrDataEmpty          = NewHTTPErrorMsg("data is empty", http.StatusBadRequest)
	ErrDataNotFound       = NewHTTPErrorMsg("data not found", http.StatusNotFound)
)

type (
	ErrStatusCode interface{ ErrStatusCode() int }   // HTTP Code
	ErrMessage    interface{ ErrMessage() string }   // user message
	ErrCode       interface{ ErrCode() int }         // error code
	ErrData       interface{ ErrData() interface{} } // other data
)

type HTTPError struct {
	Message    string      `json:"message"`
	Code       int         `json:"code"`
	Data       interface{} `json:"data,omitempty"`
	StatusCode int         `json:"-"`
	Err        error       `json:"-"`
}

func (s *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError: %#v", s)
}

func (s *HTTPError) ErrStatusCode() int {
	return s.StatusCode
}

func (s *HTTPError) ErrMessage() string {
	return s.Message
}

func (s *HTTPError) ErrCode() int {
	return s.Code
}

func NewHTTPError(err error) *HTTPError {
	if hErr, ok := err.(*HTTPError); ok {
		return hErr
	}

	standErr := &HTTPError{
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
		Err:        err,
	}
	if err2, ok := err.(ErrStatusCode); ok {
		standErr.StatusCode = err2.ErrStatusCode()
	}
	if err2, ok := err.(ErrMessage); ok {
		standErr.Message = err2.ErrMessage()
	}
	if err2, ok := err.(ErrCode); ok {
		standErr.Code = err2.ErrCode()
	}
	if err2, ok := err.(ErrData); ok {
		standErr.Data = err2.ErrData()
	}
	return standErr
}

func NewHTTPErrorStatusCode(err error, status int) *HTTPError {
	hErr := NewHTTPError(err)
	hErr.StatusCode = status
	return hErr
}

func NewHTTPError200(err error) *HTTPError {
	return NewHTTPErrorStatusCode(err, http.StatusOK)
}

func NewHTTPErrorMsg(msg string, status int) *HTTPError {
	return NewHTTPErrorStatusCode(errors.New(msg), status)
}
