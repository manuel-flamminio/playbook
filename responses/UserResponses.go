package responses

import (
	"playbook/entities"

	"github.com/google/uuid"
)

type BaseUserInfoResponse struct {
	ID          uuid.UUID `json:"id" example:"01976595-4044-7426-b6fa-f64173211b94"`
	DisplayName string    `json:"display_name,omitempty" example:"The Lion"`
	Username    string    `json:"username,omitempty" example:"the.lion@example.com"`
}

type ElasticSearchUserListResponse struct {
	Total int
	Page  int
	Users []*BaseUserInfoResponse
}

func NewBaseUserInfoResponse(user *entities.User) *BaseUserInfoResponse {
	return &BaseUserInfoResponse{
		ID:          user.ID,
		DisplayName: user.DisplayName,
		Username:    user.Username,
	}
}
