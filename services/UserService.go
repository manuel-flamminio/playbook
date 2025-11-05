package services

import (
	"errors"
	"playbook/entities"
	"playbook/filters"
	"playbook/mappers"
	"playbook/repositories"
	"playbook/requests"
	"playbook/responses"
	"playbook/utils"

	"github.com/go-crypt/crypt"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository       repositories.UserRepositoryInterface
	dtoMapper            mappers.DtoMapperInterface
	elasticSearchWrapper ElasticSearchWrapperInterface
}

func (u *UserService) Create(user *requests.CreateUserRequest) (*entities.User, error) {
	if !utils.IsUsernameValid(user.Username) {
		return nil, errors.New("Username is not valid")
	}
	if user.DisplayName == "" {
		user.DisplayName = user.Username
	}
	userToAdd := u.dtoMapper.UserCreateToEntity(user)
	err := u.userRepository.Create(userToAdd)
	if err != nil {
		return nil, err
	}

	if err := u.elasticSearchWrapper.IndexUser(userToAdd); err != nil {
		errDelete := u.userRepository.Delete(userToAdd)
		if errDelete != nil {
			return nil, errDelete
		}
		return nil, err
	}
	return userToAdd, nil
}

func (u *UserService) UpdatePassword(userUUID uuid.UUID, userDto *requests.UpdateUserPasswordRequest) (*entities.User, error) {
	userToUpdate, err := u.userRepository.GetById(userUUID)
	if err != nil {
		return nil, err
	}

	userToUpdate.Password = userDto.Password
	return userToUpdate, u.userRepository.Update(userToUpdate)
}

func (u *UserService) UpdateDisplayName(userUUID uuid.UUID, userDto *requests.UpdateUserDisplayNameRequest) (*entities.User, error) {
	userToUpdate, err := u.userRepository.GetById(userUUID)
	if err != nil {
		return nil, err
	}

	userToUpdate.DisplayName = userDto.DisplayName
	return userToUpdate, u.userRepository.Update(userToUpdate)
}

func (u *UserService) GetByUsername(username string) (*entities.User, error) {
	return u.userRepository.GetByUsername(username)
}

func (u *UserService) SearchUser(filters *filters.UserQueryFilters) (*responses.ElasticSearchUserListResponse, error) {
	return u.elasticSearchWrapper.SearchUsers(filters)
}

func (u *UserService) Delete(user *entities.User) error {
	if err := u.userRepository.Delete(user); err != nil {
		return err
	}
	return u.elasticSearchWrapper.DeleteUser(user.ID)
}

func CreateUserService(userRepository repositories.UserRepositoryInterface, dtoMapper mappers.DtoMapperInterface, elasticSearchWrapper ElasticSearchWrapperInterface) *UserService {
	return &UserService{userRepository: userRepository, dtoMapper: dtoMapper, elasticSearchWrapper: elasticSearchWrapper}
}

// Login godoc
// @Summary      Get Auth token
// @Description  Get Auth token
// @Tags         Auth
// @Accept       json
// @Produce      json
//
// @Param request body requests.UserLoginRequest true "User Credentials"
//
// @Success      200  {object}  responses.BaseUserInfoResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /login [post]
func (u *UserService) Login(user *entities.User) (*entities.User, error) {
	foundUser, err := u.userRepository.GetByUsername(user.Username)
	if err != nil {
		return nil, err
	}

	isPasswordValid, err := crypt.CheckPassword(user.Password, foundUser.Password)
	if err != nil {
		return nil, err
	}

	if !isPasswordValid {
		return nil, errors.New("Invalid Credentials")
	}

	return foundUser, nil
}
