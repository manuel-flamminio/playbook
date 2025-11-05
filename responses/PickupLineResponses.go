package responses

import (
	"playbook/entities"
	"time"

	"github.com/google/uuid"
)

type SinglePickupLineInfoResponse struct {
	ID           uuid.UUID             `json:"id" example:"cc4e0e71-d67b-4249-8a64-e75196835df4"`
	Title        string                `json:"title,omitempty" example:"The Baker"`
	Content      string                `json:"content,omitempty" example:"Yo, are you a baker? Because you are a cutie pie"`
	Tags         []*entities.Tag       `json:"tags,omitempty"`
	Statistics   *entities.Statistic   `json:"statistics"`
	UserReaction *entities.Reaction    `json:"reaction"`
	Visible      bool                  `json:"visible"`
	UpdatedAt    time.Time             `json:"updated_at,omitempty" example:"2025-01-01T11:23:00Z"`
	User         *BaseUserInfoResponse `json:"user"`
}

type ElasticSearchPickupLineResponse struct {
	Total int
	Users []*SinglePickupLineInfoResponse
}

func NewSinglePickupLineInfoResponse(pickupLine *entities.PickupLine) *SinglePickupLineInfoResponse {
	return &SinglePickupLineInfoResponse{
		ID:           pickupLine.ID,
		Title:        pickupLine.Title,
		Content:      pickupLine.Content,
		Tags:         pickupLine.Tags,
		Statistics:   pickupLine.Statistics,
		UserReaction: pickupLine.UserReaction,
		Visible:      pickupLine.Visible,
		UpdatedAt:    pickupLine.UpdatedAt,
		User:         NewBaseUserInfoResponse(pickupLine.User),
	}
}
