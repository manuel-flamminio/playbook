package elastic_dtos

import (
	"time"

	"github.com/google/uuid"
)

type PickupLineElasticDTO struct {
	Id                uuid.UUID   `json:"id"`
	Title             string      `json:"title"`
	Content           string      `json:"content"`
	Tags              []uuid.UUID `json:"tags"`
	UserId            uuid.UUID   `json:"userId"`
	Username          string      `json:"username"`
	DisplayName       string      `json:"display_name"`
	Visible           bool        `json:"visible"`
	Starred           bool        `json:"starred"`
	NumberOfSuccesses int         `json:"numberOfSuccesses"`
	NumberOfFailures  int         `json:"numberOfFailures"`
	NumberOfTries     int         `json:"numberOfTries"`
	SuccessPercentage float64     `json:"successPercentage"`
	UpdatedAt         time.Time   `json:"updatedAt"`
}

type UpdatePickupLineElasticDTO struct {
	Title     string      `json:"title"`
	Content   string      `json:"content"`
	Tags      []uuid.UUID `json:"tags"`
	Visible   bool        `json:"visible"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

type TagElasticDTO struct {
	Id     uuid.UUID `json:"id"`
	Name   string    `json:"name"`
	UserId uuid.UUID `json:"userId"`
}

type UserElasticDTO struct {
	Id          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"display_name"`
}
