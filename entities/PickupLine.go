package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Visibility string

const (
	All        Visibility = "ALL"
	Visible    Visibility = "VISIBLE"
	NotVisible Visibility = "NOT_VISIBLE"
)

type PickupLine struct {
	ID           uuid.UUID   `gorm:"type:uuid;primary_key;" json:"id"`
	Title        string      `gorm:"type:varchar(255);" json:"title,omitempty"`
	Content      string      `gorm:"type:text;" json:"content,omitempty"`
	Tags         []*Tag      `gorm:"many2many:pickup_line_tags;constraint:OnDelete:CASCADE" json:"tags,omitempty"`
	Statistics   *Statistic  `gorm:"-" json:"statistics"`
	Reactions    []*Reaction `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
	UserReaction *Reaction   `json:"reaction"`
	Visible      bool        `gorm:"not null" json:"visible"`
	UpdatedAt    time.Time   `json:"updated_at,omitempty"`
	User         *User       `gorm:"foreignkey:user_id"`
	UserID       uuid.UUID
}

func (pickupLine *PickupLine) BeforeCreate(tx *gorm.DB) error {
	if pickupLine.ID != uuid.Nil {
		return nil
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	pickupLine.ID = uuid

	return nil
}

func (u *PickupLine) IsEqualTo(otherPickupLine *PickupLine) bool {
	return u.ID == otherPickupLine.ID &&
		u.Title == otherPickupLine.Title &&
		u.Content == otherPickupLine.Content &&
		u.Visible == otherPickupLine.Visible &&
		u.hasSameTags(otherPickupLine)
}

func (u *PickupLine) hasSameTags(otherPickupLine *PickupLine) bool {
	if len(u.Tags) != len(otherPickupLine.Tags) {
		return false
	}

	tagMap := make(map[string]bool, len(u.Tags))
	for _, tag := range u.Tags {
		tagMap[tag.ID.String()] = true
	}

	for _, tag := range otherPickupLine.Tags {
		_, ok := tagMap[tag.ID.String()]
		if !ok {
			return false
		}
	}

	return true
}
