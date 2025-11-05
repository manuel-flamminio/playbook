package services

import (
	"errors"
	"playbook/entities"
	"playbook/filters"
	"playbook/mocks"
	"playbook/requests"
	"playbook/responses"
	"playbook/utils"
	"testing"

	"go.uber.org/mock/gomock"
)

func getRandomCreateUserRequest() *requests.CreateUserRequest {
	return &requests.CreateUserRequest{
		DisplayName: utils.GetRandomString(10),
		Username:    utils.GetRandomUsername(),
		Password:    utils.GetRandomString(18),
	}
}

func TestAddUserSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUser := getRandomCreateUserRequest()
	userEntity := &entities.User{
		DisplayName: randomUser.DisplayName,
		Username:    randomUser.Username,
		Password:    randomUser.Password,
	}
	expectedCreatedUser := &entities.User{
		DisplayName: randomUser.DisplayName,
		Username:    randomUser.Username,
		Password:    randomUser.Password,
	}

	mockDtoMapper.EXPECT().UserCreateToEntity(randomUser).Return(userEntity)
	mockRepository.EXPECT().Create(userEntity).Return(nil)
	mockElasticSearchWrapper.EXPECT().IndexUser(userEntity).Return(nil)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if createdUser, err := service.Create(randomUser); err != nil || !createdUser.IsEqualTo(expectedCreatedUser) {
		t.Errorf("Error occurred in user creation service")
		return
	}
}

func TestAddUserFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUser := getRandomCreateUserRequest()
	userEntity := &entities.User{
		DisplayName: randomUser.DisplayName,
		Username:    randomUser.Username,
		Password:    randomUser.Password,
	}
	expectedError := errors.New("error")

	mockDtoMapper.EXPECT().UserCreateToEntity(randomUser).Return(userEntity)
	mockRepository.EXPECT().Create(userEntity).Return(expectedError)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if _, err := service.Create(randomUser); err == nil || err != expectedError {
		t.Errorf("Error occurred in user creation service")
		return
	}
}

func TestAddUserElasticFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUser := getRandomCreateUserRequest()
	userEntity := &entities.User{
		DisplayName: randomUser.DisplayName,
		Username:    randomUser.Username,
		Password:    randomUser.Password,
	}
	expectedError := errors.New("error")

	mockDtoMapper.EXPECT().UserCreateToEntity(randomUser).Return(userEntity)
	mockRepository.EXPECT().Create(userEntity).Return(nil)
	mockElasticSearchWrapper.EXPECT().IndexUser(userEntity).Return(expectedError)
	mockRepository.EXPECT().Delete(userEntity).Return(nil)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if _, err := service.Create(randomUser); err == nil || err != expectedError {
		t.Errorf("Error occurred in user creation service")
		return
	}
}

func TestAddUserElasticFailureAndDatabaseFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUser := getRandomCreateUserRequest()
	userEntity := &entities.User{
		DisplayName: randomUser.DisplayName,
		Username:    randomUser.Username,
		Password:    randomUser.Password,
	}
	expectedError := errors.New("error")
	otherError := errors.New("error")

	mockDtoMapper.EXPECT().UserCreateToEntity(randomUser).Return(userEntity)
	mockRepository.EXPECT().Create(userEntity).Return(nil)
	mockElasticSearchWrapper.EXPECT().IndexUser(userEntity).Return(otherError)
	mockRepository.EXPECT().Delete(userEntity).Return(expectedError)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if _, err := service.Create(randomUser); err == nil || err != expectedError {
		t.Errorf("Error occurred in user creation service")
		return
	}
}

func TestAddUserMissingDisplayNameSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUser := getRandomCreateUserRequest()
	randomUser.DisplayName = ""
	userEntity := &entities.User{
		DisplayName: randomUser.Username,
		Username:    randomUser.Username,
		Password:    randomUser.Password,
	}
	expectedCreatedUser := &entities.User{
		DisplayName: randomUser.Username,
		Username:    randomUser.Username,
		Password:    randomUser.Password,
	}

	mockDtoMapper.EXPECT().UserCreateToEntity(randomUser).Return(userEntity)
	mockRepository.EXPECT().Create(userEntity).Return(nil)
	mockElasticSearchWrapper.EXPECT().IndexUser(userEntity).Return(nil)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if createdUser, err := service.Create(randomUser); err != nil || !createdUser.IsEqualTo(expectedCreatedUser) {
		t.Errorf("Error occurred in user creation service")
		return
	}
}

func TestAddUserInvalidUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUser := getRandomCreateUserRequest()
	randomUser.Username = "invalidUsername  aa"
	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if user, err := service.Create(randomUser); err == nil || user != nil {
		t.Errorf("Expected invalid username error")
		return
	}
}

func TestAddUserBlankUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUser := getRandomCreateUserRequest()
	randomUser.Username = "    "
	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if _, err := service.Create(randomUser); err == nil {
		t.Errorf("Expected invalid username error")
		return
	}
}

func TestUpdateDisplayName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUpdateDisplayNameRequest := &requests.UpdateUserDisplayNameRequest{DisplayName: utils.GetRandomString(10)}
	randomUser := utils.GetRandomValidUser()

	mockRepository.EXPECT().GetById(randomUser.ID).Return(randomUser, nil)
	randomUser.DisplayName = randomUpdateDisplayNameRequest.DisplayName
	mockRepository.EXPECT().Update(randomUser).Return(nil)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if _, err := service.UpdateDisplayName(randomUser.ID, randomUpdateDisplayNameRequest); err != nil {
		t.Errorf("Error occurred while updating user DisplayName")
		return
	}
}

func TestUpdateDisplayNameRetrieveFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUpdateDisplayNameRequest := &requests.UpdateUserDisplayNameRequest{DisplayName: utils.GetRandomString(10)}
	randomUser := utils.GetRandomValidUser()
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetById(randomUser.ID).Return(nil, expectedError)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if user, err := service.UpdateDisplayName(randomUser.ID, randomUpdateDisplayNameRequest); err == nil || user != nil {
		t.Errorf("Error occurred while updating user DisplayName")
		return
	}
}

func TestUpdatePassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUpdatePasswordRequest := &requests.UpdateUserPasswordRequest{Password: utils.GetRandomString(18)}
	randomUser := utils.GetRandomValidUser()

	mockRepository.EXPECT().GetById(randomUser.ID).Return(randomUser, nil)
	randomUser.Password = randomUpdatePasswordRequest.Password
	mockRepository.EXPECT().Update(randomUser).Return(nil)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if _, err := service.UpdatePassword(randomUser.ID, randomUpdatePasswordRequest); err != nil {
		t.Errorf("Error occurred while updating user Password")
		return
	}
}

func TestUpdatePasswordRetrieveFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUpdatePasswordRequest := &requests.UpdateUserPasswordRequest{Password: utils.GetRandomString(18)}
	randomUser := utils.GetRandomValidUser()
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetById(randomUser.ID).Return(nil, expectedError)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if user, err := service.UpdatePassword(randomUser.ID, randomUpdatePasswordRequest); err == nil || user != nil {
		t.Errorf("Error occurred while updating user Password")
		return
	}
}

func TestGetByUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	username := utils.GetRandomUsername()
	mockRepository.EXPECT().GetByUsername(username).Return(nil, nil)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if _, err := service.GetByUsername(username); err != nil {
		t.Errorf("Error occurred while getting user by username")
		return
	}
}

func TestSearchUsers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	userList := &responses.ElasticSearchUserListResponse{}

	filters := &filters.UserQueryFilters{}
	mockElasticSearchWrapper.EXPECT().SearchUsers(filters).Return(userList, nil)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if _, err := service.SearchUser(filters); err != nil {
		t.Errorf("Error occurred while deleting user")
		return
	}
}

func TestDelete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	user := utils.GetRandomValidUser()
	mockRepository.EXPECT().Delete(user).Return(nil)
	mockElasticSearchWrapper.EXPECT().DeleteUser(user.ID)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if err := service.Delete(user); err != nil {
		t.Errorf("Error occurred while deleting user")
		return
	}
}

func TestDeleteFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	expectedError := errors.New("error")
	user := utils.GetRandomValidUser()
	mockRepository.EXPECT().Delete(user).Return(expectedError)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if err := service.Delete(user); err == nil || err != expectedError {
		t.Errorf("Error occurred while deleting user")
		return
	}
}

func TestLoginSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUser := getRandomCreateUserRequest()
	userEntity := &entities.User{
		DisplayName: randomUser.DisplayName,
		Username:    randomUser.Username,
		Password:    randomUser.Password,
	}
	foundUser := &entities.User{
		DisplayName: randomUser.DisplayName,
		Username:    randomUser.Username,
		Password:    randomUser.Password,
	}
	foundUser.HashPassword()

	mockRepository.EXPECT().GetByUsername(userEntity.Username).Return(foundUser, nil)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if _, err := service.Login(userEntity); err != nil {
		t.Errorf("Error occurred in user login")
		return
	}
}

func TestLoginPasswordNotValid(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUser := getRandomCreateUserRequest()
	userEntity := &entities.User{
		DisplayName: randomUser.DisplayName,
		Username:    randomUser.Username,
		Password:    randomUser.Password,
	}
	foundUser := &entities.User{
		DisplayName: randomUser.DisplayName,
		Username:    randomUser.Username,
		Password:    utils.GetRandomString(20),
	}
	foundUser.HashPassword()

	mockRepository.EXPECT().GetByUsername(userEntity.Username).Return(foundUser, nil)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if _, err := service.Login(userEntity); err == nil {
		t.Errorf("Error occurred in user login")
		return
	}
}

func TestLoginRetrieveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockUserRepositoryInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockElasticSearchWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)

	randomUser := getRandomCreateUserRequest()
	userEntity := &entities.User{
		DisplayName: randomUser.DisplayName,
		Username:    randomUser.Username,
		Password:    randomUser.Password,
	}
	foundUser := &entities.User{
		DisplayName: randomUser.DisplayName,
		Username:    randomUser.Username,
		Password:    utils.GetRandomString(20),
	}
	foundUser.HashPassword()
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetByUsername(userEntity.Username).Return(nil, expectedError)

	service := CreateUserService(mockRepository, mockDtoMapper, mockElasticSearchWrapper)
	if _, err := service.Login(userEntity); err == nil {
		t.Errorf("Error occurred in user login")
		return
	}
}
