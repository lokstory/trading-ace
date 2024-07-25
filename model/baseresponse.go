package model

type BaseResponse struct {
	Status Status                 `json:"status"`
	Extra  map[string]interface{} `json:"extra,omitempty"`
}

func NewBaseResponse(status Status) *BaseResponse {
	return &BaseResponse{
		Status: status,
	}
}
