package mappers

import (
	"playbook/entities"
	"playbook/requests"
	"playbook/utils"
	"testing"
)

func TestUserCreateToUserEntity(t *testing.T) {
	userCreate := &requests.CreateUserRequest{
		DisplayName:      utils.GetRandomString(10),
		Username:         utils.GetRandomUsername(),
		EncodedUserImage: utils.GetRandomString(13),
		Password:         utils.GetRandomString(18),
	}
	expectedUserEntity := &entities.User{
		Username:    userCreate.Username,
		DisplayName: userCreate.DisplayName,
		Password:    userCreate.Password,
	}

	mapper := GetNewDtoMapper()
	if mappedEntity := mapper.UserCreateToEntity(userCreate); !mappedEntity.IsEqualTo(expectedUserEntity) || mappedEntity.Password != expectedUserEntity.Password {
		t.Errorf("Error occurred in dto mapper")
		return
	}
}

func TestTagBodyRequestToEntity(t *testing.T) {
	userUUID, _ := utils.GetRandomUUID()
	tagCreate := &requests.TagBodyRequest{
		Name:        utils.GetRandomTagName(),
		Description: utils.GetRandomTagDescription(),
	}
	expectedTagEntity := &entities.Tag{
		Name:        tagCreate.Name,
		Description: tagCreate.Description,
		UserId:      userUUID,
	}

	mapper := GetNewDtoMapper()
	if mappedEntity := mapper.TagBodyRequestToEntity(tagCreate, userUUID); !mappedEntity.IsEqualTo(expectedTagEntity) {
		t.Errorf("Error occurred in dto mapper")
		return
	}
}

func TestGetPickupLineEntityFromDTO(t *testing.T) {
	pickupLineCreate := &requests.PickupLineBodyRequest{
		Title:   utils.GetRandomString(10),
		Content: utils.GetRandomString(20),
		Visible: utils.GetRandomBool(),
	}
	expectedPickupLineEntity := &entities.PickupLine{
		Title:   pickupLineCreate.Title,
		Content: pickupLineCreate.Content,
		Visible: pickupLineCreate.Visible,
	}

	mapper := GetNewDtoMapper()
	if mappedEntity := mapper.GetPickupLineEntityFromDTO(pickupLineCreate); !mappedEntity.IsEqualTo(expectedPickupLineEntity) {
		t.Errorf("Error occurred in dto mapper")
		return
	}
}
