package mappers

import (
	"playbook/entities"
	"playbook/requests"

	"github.com/google/uuid"
)

type DtoMapperInterface interface {
	UserCreateToEntity(*requests.CreateUserRequest) *entities.User
	TagBodyRequestToEntity(tagRequest *requests.TagBodyRequest, userUUID uuid.UUID) *entities.Tag
	GetPickupLineEntityFromDTO(pickupLineDTO *requests.PickupLineBodyRequest) *entities.PickupLine
}

type DtoMapper struct{}

func (d *DtoMapper) UserCreateToEntity(user *requests.CreateUserRequest) *entities.User {
	return &entities.User{
		DisplayName: user.DisplayName,
		Username:    user.Username,
		Password:    user.Password,
	}
}

func (d *DtoMapper) TagBodyRequestToEntity(tagRequest *requests.TagBodyRequest, userUUID uuid.UUID) *entities.Tag {
	return &entities.Tag{
		Name:        tagRequest.Name,
		Description: tagRequest.Description,
		UserId:      userUUID,
	}
}

func (d *DtoMapper) GetPickupLineEntityFromDTO(pickupLineDTO *requests.PickupLineBodyRequest) *entities.PickupLine {
	return &entities.PickupLine{
		Title:   pickupLineDTO.Title,
		Content: pickupLineDTO.Content,
		Visible: pickupLineDTO.Visible,
	}
}

func GetNewDtoMapper() DtoMapperInterface {
	return &DtoMapper{}
}
