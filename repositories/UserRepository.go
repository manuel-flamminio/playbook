package repositories

import (
	"playbook/database"
	"playbook/entities"

	"github.com/google/uuid"
)

type UserRepositoryInterface interface {
	Create(user *entities.User) error
	Update(user *entities.User) error
	GetByUsername(username string) (*entities.User, error)
	GetById(userUUID uuid.UUID) (*entities.User, error)
	Delete(user *entities.User) error
}

type UserRepository struct {
	database *database.Database
}

type UserRepositoryMock struct {
}

func (u *UserRepository) Create(user *entities.User) error {
	return u.database.CreateUser(user)
}

func (u *UserRepository) Update(user *entities.User) error {
	return u.database.UpdateUser(user)
}

func CreateUserRepository(database *database.Database) *UserRepository {
	return &UserRepository{database: database}
}

func (u *UserRepository) GetByUsername(username string) (*entities.User, error) {
	return u.database.GetUserByUsername(username)
}

func (u *UserRepository) GetById(userUUID uuid.UUID) (*entities.User, error) {
	return u.database.GetUserById(userUUID)
}

func (u *UserRepository) Delete(user *entities.User) error {
	return u.database.DeleteUser(user)
}
