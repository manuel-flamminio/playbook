package requests

import (
	"playbook/entities"
	"playbook/filters"
)

type PickupLineBodyRequest struct {
	Title   string          `json:"title,omitempty" example:"The Baker"`
	Content string          `json:"content,omitempty" example:"Yo, are you a baker? Because you are a cutie pie"`
	Tags    []*TagIdRequest `json:"tags,omitempty"`
	Visible bool            `json:"visible" example:"true"`
}

type PickupLineFilters struct {
	Page              int                 `form:"page"`
	Title             string              `form:"title"`
	Starred           bool                `form:"starred"`
	Visible           entities.Visibility `form:"visibility"`
	OnlyUpvoted       bool                `form:"only_upvoted" description:"this filter works only for the currently requesting user"`
	Content           string              `form:"content"`
	Tags              []string            `form:"tags[]"`
	SuccessPercentage float64             `form:"success_percentage"`
	UserId            string              `form:"user_id"`
	SortingType       filters.SortingType `form:"sorting_type"`
}
