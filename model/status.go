package model

type Status string

const (
	Success        Status = "SUCCESS"
	ParameterError        = "PARAMETER_ERROR"
	ServerError           = "SERVER_ERROR"
	OperateFailed         = "OPERATE_FAILED"
	NotFound              = "NOT_FOUND"
)
