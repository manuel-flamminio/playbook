package filters

import (
	"playbook/entities"

	"github.com/google/uuid"
)

type SortingType string

const (
	New           SortingType = "NEW"
	BestOfAllTime SortingType = "BEST_OF_ALL_TIME"
	Trending      SortingType = "TRENDING"
	Random        SortingType = "RANDOM"
)

type PickupLineQueryFilters struct {
	Page              int
	UserId            uuid.UUID
	Title             string
	Starred           bool
	OnlyUpvoted       bool
	Visible           entities.Visibility
	Tags              []string
	Content           string
	SuccessPercentage float64
	SortingType       SortingType
}
