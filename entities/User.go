package entities

import (
	"playbook/security"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID          uuid.UUID    `gorm:"type:uuid;primary_key;" json:"id"`
	DisplayName string       `gorm:"type:varchar(255)" json:"display_name,omitempty"`
	Username    string       `gorm:"type:varchar(255);uniqueIndex" json:"username,omitempty"`
	Password    string       `gorm:"type:varchar(255);" json:"password,omitempty"`
	Tags        []Tag        `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"tags,omitempty"`
	PickupLines []PickupLine `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"pickup_lines,omitempty"`
	Reaction    []*Reaction  `gorm:"constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"-"`
}

func (user *User) BeforeSave(tx *gorm.DB) error {
	if user.ID != uuid.Nil {
		return nil
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	user.ID = uuid

	return nil
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	err := user.HashPassword()
	if err != nil {
		return err
	}

	return nil
}

func (u *User) HashPassword() error {
	hash, err := security.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hash

	return nil
}

func (u *User) IsEqualTo(otherUser *User) bool {
	return u.ID == otherUser.ID &&
		u.DisplayName == otherUser.DisplayName &&
		u.Username == otherUser.Username
}
