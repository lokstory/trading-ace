package model

type OKResponse struct {
	*BaseResponse
	Data interface{} `json:"data"`
}

func NewOKResponse(data interface{}) *OKResponse {
	return &OKResponse{
		BaseResponse: NewBaseResponse(Success),
		Data:         data,
	}
}
