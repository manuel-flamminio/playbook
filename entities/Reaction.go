package entities

import "github.com/google/uuid"

type Vote string

const (
	None     Vote = "NONE"
	Upvote   Vote = "UPVOTE"
	Downvote Vote = "DOWNVOTE"
)

type Reaction struct {
	PickupLineId uuid.UUID `gorm:"type:uuid;primary_key;" json:"-"`
	UserId       uuid.UUID `gorm:"type:uuid;primary_key;" json:"-"`
	Starred      bool      `gorm:"not null" json:"starred"`
	Vote         Vote      `gorm:"type:varchar(100);default:'NONE';not null" json:"vote"`
}
