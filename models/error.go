package models

import "net/http"

type Errors struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

func (e *Errors) Error() string {
	return e.Message
}

func NewErrors(code int, msg string, data interface{}) *Errors {
	return &Errors{
		Code:    code,
		Message: msg,
		Errors:  data,
	}
}

func NewBadRequest(data interface{}) *Errors {
	return &Errors{
		Code:    http.StatusBadRequest,
		Message: "bad request",
		Errors:  data,
	}
}

func NewBadRequestParameterID() *Errors {
	return &Errors{
		Code:    http.StatusBadRequest,
		Message: "bad request",
		Errors:  "error parameter id",
	}
}

func NewFailedToConvertData() *Errors {
	return &Errors{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
		Errors:  "cannot convert id param",
	}
}

func NewInternalServerError(data interface{}) *Errors {
	return &Errors{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
		Errors:  data,
	}
}

func NewNotFound(data interface{}) *Errors {
	return &Errors{
		Code:    http.StatusNotFound,
		Message: "not found",
		Errors:  data,
	}
}
