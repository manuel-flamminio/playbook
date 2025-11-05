package responses

type GenericSuccessResponse struct {
	Success bool `json:"success" example:"true"`
}

func NewGenericSuccessResponse() *GenericSuccessResponse {
	return &GenericSuccessResponse{Success: true}
}
