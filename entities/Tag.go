package entities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tag struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;" json:"id" example:"01976595-4044-7426-b6fa-f64173211b94"`
	Name            string    `gorm:"type:varchar(255);" json:"name,omitempty" example:"Risky"`
	Description     string    `gorm:"type:text;" json:"description,omitempty" example:"PickupLine that could send you to jail"`
	UserId          uuid.UUID `json:"user_id,omitempty" example:"01976595-4044-7426-b6fa-f64173211b94"`
	ElasticSearchId string    `gorm:"type:varchar(255);" json:"-"`
}

func (tag *Tag) BeforeCreate(tx *gorm.DB) error {
	if tag.ID != uuid.Nil {
		return nil
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	tag.ID = uuid

	return nil
}

func (t *Tag) IsEqualTo(otherTag *Tag) bool {
	return t.ID == otherTag.ID &&
		t.Name == otherTag.Name &&
		t.Description == otherTag.Description
}
