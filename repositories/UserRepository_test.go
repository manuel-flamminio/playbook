package repositories

import (
	"playbook/database"
	"playbook/entities"
	"playbook/utils"
	"testing"

	"github.com/google/uuid"
)

func getRandomTestUser() (*entities.User, error) {
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	return &entities.User{
		ID:       uuid,
		Username: utils.GetRandomUsername(),
		Password: "testeroni",
	}, nil
}

func TestSelectUser(t *testing.T) {
	user, err := getRandomTestUser()
	if err != nil {
		t.Errorf("Unable to create random user")
		return
	}

	rawDb, cleanup := database.InitializeAndGetTestDatabase()
	defer cleanup()
	tx := rawDb.Begin()

	result := tx.Create(user)
	if result.Error != nil {
		t.Errorf("Unable to insert random user inside database")
		return
	}

	db := database.GetTestDatabaseConnection(tx)
	repository := CreateUserRepository(db)
	selectedUserById, err := repository.GetById(user.ID)
	if err != nil {
		t.Errorf("Cannot Select User")
		tx.Rollback()
		return
	}

	if selectedUserById.ID != user.ID {
		t.Errorf("Got wrong user")
	}

	selectedUserByUsername, err := repository.GetByUsername(user.Username)
	if err != nil {
		t.Errorf("Cannot Select User")
		tx.Rollback()
		return
	}

	if selectedUserByUsername.Username != user.Username || selectedUserById.Username != selectedUserByUsername.Username {
		t.Errorf("Got wrong user by username")
	}

	tx.Rollback()
}

func TestAddUser(t *testing.T) {
	user, err := getRandomTestUser()
	if err != nil {
		t.Errorf("Unable to create random user")
		return
	}

	rawDb, cleanup := database.InitializeAndGetTestDatabase()
	defer cleanup()
	tx := rawDb.Begin()

	user.ID = uuid.Nil
	db := database.GetTestDatabaseConnection(tx)
	repository := CreateUserRepository(db)
	result := repository.Create(user)
	if result != nil {
		t.Errorf("Unable to insert random user inside database")
		return
	}

	selectedUser := &entities.User{}
	tx.First(&selectedUser, "id = ?", user.ID)
	if selectedUser.ID == uuid.Nil {
		t.Errorf("Selected user uuid is not valid")
	}

	if selectedUser.ID != user.ID {
		t.Errorf("Selected user is not the inserted one")
	}

	if selectedUser.Password != user.Password {
		t.Errorf("Selected user password is not correct")
	}

	tx.Rollback()
}

func TestUpdateUser(t *testing.T) {
	user, err := getRandomTestUser()
	if err != nil {
		t.Errorf("Unable to create random user")
		return
	}

	rawDb, cleanup := database.InitializeAndGetTestDatabase()
	defer cleanup()
	tx := rawDb.Begin()

	user.ID = uuid.Nil
	db := database.GetTestDatabaseConnection(tx)
	repository := CreateUserRepository(db)
	result := repository.Create(user)
	if result != nil {
		t.Errorf("Unable to insert random user inside database")
		return
	}

	oldUsername := user.Username
	user.Username = "nerwest"
	result = repository.Update(user)
	if result != nil {
		t.Errorf("Unable to update random user inside database")
		return
	}

	selectedUser := &entities.User{}
	tx.First(&selectedUser, "id = ?", user.ID)
	if selectedUser.ID != user.ID {
		t.Errorf("Selected user is not the inserted one")
	}

	if selectedUser.Username == oldUsername || selectedUser.Username != user.Username {
		t.Errorf("Selected user username was not updated correctly")
	}

	tx.Rollback()
}
