package model

import (
	"net/http"
)

type AppError struct {
	error
	*BaseResponse
	Data           interface{} `json:"data"`
	message        string
	err            error
	HTTPStatusCode int `json:"-"`
}

func (e *AppError) Err(err error) *AppError {
	if err != nil {
		e.err = err
	}

	return e
}

func (e *AppError) Message(message string) *AppError {
	e.message = message

	return e
}

func (e *AppError) StatusCode(statusCode int) *AppError {
	e.HTTPStatusCode = statusCode

	return e
}

func (e *AppError) Error() (text string) {
	if e.err != nil {
		text = e.err.Error()
		return
	}

	if e.message != "" {
		text = e.message
	}

	return
}

func NewAppError(status Status) *AppError {
	return &AppError{
		BaseResponse:   NewBaseResponse(status),
		HTTPStatusCode: http.StatusOK,
	}
}

func NewParameterError() *AppError {
	return &AppError{
		BaseResponse:   NewBaseResponse(ParameterError),
		HTTPStatusCode: http.StatusBadRequest,
	}
}

func NewServerError() *AppError {
	return &AppError{
		BaseResponse:   NewBaseResponse(ServerError),
		HTTPStatusCode: http.StatusInternalServerError,
	}
}

var _ error = &AppError{}
