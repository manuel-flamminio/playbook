package repositories

import (
	"errors"
	"playbook/database"
	"playbook/entities"
	"playbook/utils"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func getRandomTestTagAndUser() (*entities.Tag, *entities.User) {
	tagUuid, err := uuid.NewV7()
	if err != nil {
		panic("could not generate uuid")
	}

	user, err := getRandomTestUser()
	if err != nil {
		panic("could not generate user")
	}

	return &entities.Tag{
		ID:              tagUuid,
		Name:            utils.GetRandomTagName(),
		Description:     utils.GetRandomTagDescription(),
		UserId:          user.ID,
		ElasticSearchId: "",
	}, user
}

func getRandomTestTagWithoutUser() *entities.Tag {
	tagUuid, err := uuid.NewV7()
	if err != nil {
		panic("could not generate uuid")
	}

	return &entities.Tag{
		ID:              tagUuid,
		Name:            utils.GetRandomTagName(),
		Description:     utils.GetRandomTagDescription(),
		ElasticSearchId: "",
	}
}

func TestSelectTag(t *testing.T) {
	tag, user := getRandomTestTagAndUser()
	tx, cleanup := database.InitializeAndGetTestDatabaseTransaction()
	defer cleanup()

	resultUser := tx.Create(user)
	resultTag := tx.Create(tag)
	if resultUser.Error != nil || resultTag.Error != nil {
		t.Errorf("Unable to insert random user or tag inside database")
		return
	}

	db := database.GetTestDatabaseConnection(tx)
	repository := CreateTagRepository(db)
	selectedTagById, err := repository.GetById(tag.ID.String())
	if err != nil {
		t.Errorf("Cannot Select Tag")
		tx.Rollback()
		return
	}

	if selectedTagById.ID != tag.ID {
		t.Errorf("Got wrong tag")
	}

	selectedTagByUsername, err := repository.GetByUserId(user.ID)
	if err != nil {
		t.Errorf("Cannot Select by UserId")
		tx.Rollback()
		return
	}

	if selectedTagByUsername[0].ID != tag.ID {
		t.Errorf("Got wrong tag by user id")
	}

	tx.Rollback()
}

func TestAddTag(t *testing.T) {
	tag, user := getRandomTestTagAndUser()
	tx, cleanup := database.InitializeAndGetTestDatabaseTransaction()
	defer cleanup()

	tag.ID = uuid.Nil
	resultUser := tx.Create(user)
	db := database.GetTestDatabaseConnection(tx)
	repository := CreateTagRepository(db)
	resultTag := repository.Create(tag)
	if resultUser.Error != nil || resultTag != nil {
		t.Errorf("Unable to insert random user or tag inside database")
		return
	}

	selectedTag := &entities.Tag{}
	tx.First(&selectedTag, "id = ?", tag.ID)
	if tag.ID == uuid.Nil {
		t.Errorf("Selected tag uuid is not valid")
	}

	if selectedTag.ID != tag.ID {
		t.Errorf("Selected tag is not the inserted one")
	}

	if selectedTag.Name != tag.Name {
		t.Errorf("Selected Tag Name is not correct")
	}

	tx.Rollback()
}

func TestUpdateTag(t *testing.T) {
	tag, user := getRandomTestTagAndUser()
	tx, cleanup := database.InitializeAndGetTestDatabaseTransaction()
	defer cleanup()

	db := database.GetTestDatabaseConnection(tx)
	resultUser := tx.Create(user)
	repository := CreateTagRepository(db)
	err := repository.Create(tag)
	if err != nil || resultUser.Error != nil {
		t.Errorf("Unable to insert random tag or user inside database")
		return
	}

	tag.Name = "new Name"
	tag.Description = "new Description"
	result := repository.Update(tag)
	if result != nil {
		t.Errorf("Unable to update random tag inside database")
		return
	}

	selectedTag := &entities.Tag{}
	tx.First(&selectedTag, "id = ?", tag.ID)
	if selectedTag.ID != tag.ID {
		t.Errorf("Selected tag is not the inserted one")
	}

	if selectedTag.Name != tag.Name || selectedTag.Description != tag.Description {
		t.Errorf("Selected Tag was not updated correctly")
	}

	tx.Rollback()
}

func TestDeleteTag(t *testing.T) {
	tag, user := getRandomTestTagAndUser()
	tx, cleanup := database.InitializeAndGetTestDatabaseTransaction()
	defer cleanup()

	tag.ID = uuid.Nil
	resultUser := tx.Create(user)
	db := database.GetTestDatabaseConnection(tx)
	repository := CreateTagRepository(db)
	resultTag := repository.Create(tag)
	if resultUser.Error != nil || resultTag != nil {
		t.Errorf("Unable to insert random user or tag inside database")
		return
	}

	selectedTag := &entities.Tag{}
	tx.First(&selectedTag, "id = ?", tag.ID)
	if tag.ID == uuid.Nil {
		t.Errorf("Selected tag uuid is not valid")
	}

	if selectedTag.ID != tag.ID {
		t.Errorf("Selected tag is not the inserted one")
	}

	if selectedTag.Name != tag.Name {
		t.Errorf("Selected Tag Name is not correct")
	}

	err := repository.Delete(tag)
	if err != nil {
		t.Errorf("Could not delete tag")
	}

	removedTag := &entities.Tag{}
	err = tx.First(&removedTag, "id = ?", tag.ID).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("Tag was not correctly deleted ")
	}

	tx.Rollback()
}

func TestSelectTagListByUser(t *testing.T) {
	numberOfTags := 10
	userOne, err := getRandomTestUser()
	if err != nil {
		t.Errorf("Could not create random user")
	}
	userTwo, err := getRandomTestUser()
	if err != nil {
		t.Errorf("Could not create random user")
	}
	userList := []*entities.User{userOne, userTwo}

	tx, cleanup := database.InitializeAndGetTestDatabaseTransaction()
	defer cleanup()
	errUserOne := tx.Create(userOne).Error
	errUserTwo := tx.Create(userTwo).Error
	if errUserOne != nil || errUserTwo != nil {
		t.Errorf("Could not insert random Users")
	}

	userToListMap := make(map[string]*[][]uuid.UUID, 0)
	for i := 0; i < 2; i++ {
		userUuidToTagUuids := make([][]uuid.UUID, 0, numberOfTags)
		activeUser := userList[i]
		userToListMap[activeUser.ID.String()] = &userUuidToTagUuids
		for j := 0; j < numberOfTags; j++ {
			tag := getRandomTestTagWithoutUser()
			tag.UserId = activeUser.ID
			userUuidToTagUuids = append(userUuidToTagUuids, []uuid.UUID{tag.ID, activeUser.ID})
			if tx.Create(tag).Error != nil {
				t.Errorf("Could Not create Tag")
			}
		}
	}

	db := database.GetTestDatabaseConnection(tx)
	repository := CreateTagRepository(db)
	for userUuid, userTags := range userToListMap {
		tagList, err := repository.GetListByTagIdAndUserIdList(*userTags)
		if err != nil {
			t.Errorf("Could Not Retrieve Tag List for user %s", userUuid)
		}

		if len(tagList) != len(*userTags) {
			t.Errorf("retrieved tagList is not complete")
		}

		tagMap := make(map[string]bool, 0)
		for _, tag := range tagList {
			tagMap[tag.ID.String()] = true
		}

		for _, tag := range *userTags {
			if _, ok := tagMap[tag[0].String()]; !ok {
				t.Errorf("Tag %s is missing", tag[0].String())
			}
		}
	}

	tx.Rollback()
}
