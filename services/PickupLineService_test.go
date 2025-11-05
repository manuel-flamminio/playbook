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

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func getRandomPickupLineBodyRequestNoTags() *requests.PickupLineBodyRequest {
	return &requests.PickupLineBodyRequest{
		Title:   utils.GetRandomString(10),
		Content: utils.GetRandomString(20),
		Tags:    []*requests.TagIdRequest{},
		Visible: utils.GetRandomBool(),
	}
}

func getRandomTagIdRequests(numberOfTags int) []*requests.TagIdRequest {
	tags := make([]*requests.TagIdRequest, 0, 3)
	for i := 0; i < numberOfTags; i++ {
		uuid, _ := utils.GetRandomUUID()
		tags = append(tags, &requests.TagIdRequest{ID: uuid.String()})
	}
	return tags
}

func TestAddPickupLineWithNoTagsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	user := &entities.User{
		ID:          uuid.Max,
		Username:    utils.GetRandomUsername(),
		DisplayName: utils.GetRandomUsername(),
	}
	randomPickupLineBodyRequest := getRandomPickupLineBodyRequestNoTags()
	expectedPickupLine := &entities.PickupLine{}
	expectedCreatedPickupLine := &entities.PickupLine{}

	mockDtoMapper.EXPECT().GetPickupLineEntityFromDTO(randomPickupLineBodyRequest).Return(expectedPickupLine)
	mockRepository.EXPECT().Create(expectedPickupLine).Return(nil)
	mockElasticWrapper.EXPECT().IndexPickupLine(expectedPickupLine).Return(nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if createdPickupLine, err := service.Create(randomPickupLineBodyRequest, user); err != nil || !createdPickupLine.IsEqualTo(expectedCreatedPickupLine) {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestAddPickupLineWithTagsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	user := &entities.User{
		ID:          uuid.Max,
		Username:    utils.GetRandomUsername(),
		DisplayName: utils.GetRandomUsername(),
	}
	randomTagsIds := getRandomTagIdRequests(3)
	tagSearchList := make([][]uuid.UUID, 0, len(randomTagsIds))
	retrievedTagList := make([]*entities.Tag, 0, len(randomTagsIds))
	for _, tagRequest := range randomTagsIds {
		tagUUID, err := uuid.Parse(tagRequest.ID)
		if err != nil {
			t.Errorf("Error occurred in PikcupLine creation service")
			return
		}

		tagSearchList = append(tagSearchList, []uuid.UUID{tagUUID, user.ID})
		retrievedTagList = append(retrievedTagList, &entities.Tag{ID: tagUUID})
	}
	randomPickupLineBodyRequest := getRandomPickupLineBodyRequestNoTags()
	randomPickupLineBodyRequest.Tags = randomTagsIds
	expectedPickupLine := &entities.PickupLine{}
	expectedCreatedPickupLine := &entities.PickupLine{
		Tags: retrievedTagList,
	}

	mockDtoMapper.EXPECT().GetPickupLineEntityFromDTO(randomPickupLineBodyRequest).Return(expectedPickupLine)
	mockTagService.EXPECT().GetListByTagIdAndUserIdList(tagSearchList).Return(retrievedTagList, nil)
	mockRepository.EXPECT().Create(expectedPickupLine).Return(nil)
	mockElasticWrapper.EXPECT().IndexPickupLine(expectedPickupLine).Return(nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if createdPickupLine, err := service.Create(randomPickupLineBodyRequest, user); err != nil || !createdPickupLine.IsEqualTo(expectedCreatedPickupLine) {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestAddPickupLineWithNoTagsCreationFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	user := &entities.User{
		ID:          uuid.Max,
		Username:    utils.GetRandomUsername(),
		DisplayName: utils.GetRandomUsername(),
	}
	randomPickupLineBodyRequest := getRandomPickupLineBodyRequestNoTags()
	expectedPickupLine := &entities.PickupLine{}
	expectedError := errors.New("could not create pickupLine")

	mockDtoMapper.EXPECT().GetPickupLineEntityFromDTO(randomPickupLineBodyRequest).Return(expectedPickupLine)
	mockRepository.EXPECT().Create(expectedPickupLine).Return(expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.Create(randomPickupLineBodyRequest, user); err == nil {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestAddPickupLineWithNoTagsElasticIndexFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	user := &entities.User{
		ID:          uuid.Max,
		Username:    utils.GetRandomUsername(),
		DisplayName: utils.GetRandomUsername(),
	}
	randomPickupLineBodyRequest := getRandomPickupLineBodyRequestNoTags()
	expectedPickupLine := &entities.PickupLine{}
	expectedError := errors.New("could not create pickupLine")

	mockDtoMapper.EXPECT().GetPickupLineEntityFromDTO(randomPickupLineBodyRequest).Return(expectedPickupLine)
	mockRepository.EXPECT().Create(expectedPickupLine).Return(nil)
	mockElasticWrapper.EXPECT().IndexPickupLine(expectedPickupLine).Return(expectedError)
	mockElasticWrapper.EXPECT().DeletePickupLine(expectedPickupLine.ID).Return(nil)
	mockRepository.EXPECT().Delete(expectedPickupLine).Return(nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.Create(randomPickupLineBodyRequest, user); err == nil || err != expectedError {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestAddPickupLineWithTagsCreationFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	user := &entities.User{
		ID:          uuid.Max,
		Username:    utils.GetRandomUsername(),
		DisplayName: utils.GetRandomUsername(),
	}
	randomTagsIds := getRandomTagIdRequests(3)
	tagSearchList := make([][]uuid.UUID, 0, len(randomTagsIds))
	retrievedTagList := make([]*entities.Tag, 0, len(randomTagsIds))
	for _, tagRequest := range randomTagsIds {
		tagUUID, err := uuid.Parse(tagRequest.ID)
		if err != nil {
			t.Errorf("Error occurred in PikcupLine creation service")
			return
		}

		tagSearchList = append(tagSearchList, []uuid.UUID{tagUUID, user.ID})
		retrievedTagList = append(retrievedTagList, &entities.Tag{ID: tagUUID})
	}
	randomPickupLineBodyRequest := getRandomPickupLineBodyRequestNoTags()
	randomPickupLineBodyRequest.Tags = randomTagsIds
	expectedPickupLine := &entities.PickupLine{}
	expectedError := errors.New("error in tags extraction")

	mockDtoMapper.EXPECT().GetPickupLineEntityFromDTO(randomPickupLineBodyRequest).Return(expectedPickupLine)
	mockTagService.EXPECT().GetListByTagIdAndUserIdList(tagSearchList).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.Create(randomPickupLineBodyRequest, user); err == nil || err != expectedError {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestDeletePickupLineSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	expectedPickupLine := &entities.PickupLine{
		ID: pickupLineUUID,
		User: &entities.User{
			ID: userUUID,
		},
	}

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(expectedPickupLine, nil)
	mockElasticWrapper.EXPECT().DeletePickupLine(expectedPickupLine.ID).Return(nil)
	mockRepository.EXPECT().Delete(expectedPickupLine).Return(nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if err := service.Delete(pickupLineUUID, userUUID); err != nil {
		t.Errorf("Error occurred in PickupLine delete service")
		return
	}
}

func TestDeletePickupLineDeleteFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	expectedPickupLine := &entities.PickupLine{
		ID: pickupLineUUID,
		User: &entities.User{
			ID: userUUID,
		},
	}
	expectedError := errors.New("could not delete pickupLine")

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(expectedPickupLine, nil)
	mockElasticWrapper.EXPECT().DeletePickupLine(expectedPickupLine.ID).Return(nil)
	mockRepository.EXPECT().Delete(expectedPickupLine).Return(expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if err := service.Delete(pickupLineUUID, userUUID); err == nil || err != expectedError {
		t.Errorf("Error occurred in PickupLine delete service")
		return
	}
}

func TestDeletePickupLineElasticSearchDeleteFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	expectedPickupLine := &entities.PickupLine{
		ID: pickupLineUUID,
		User: &entities.User{
			ID: userUUID,
		},
	}
	expectedError := errors.New("could not delete pickupLine")

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(expectedPickupLine, nil)
	mockElasticWrapper.EXPECT().DeletePickupLine(expectedPickupLine.ID).Return(expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if err := service.Delete(pickupLineUUID, userUUID); err == nil || err != expectedError {
		t.Errorf("Error occurred in PickupLine delete service")
		return
	}
}

func TestDeletePickupLineDeleteGetPickupLineFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	expectedError := errors.New("could not delete pickupLine")

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if err := service.Delete(pickupLineUUID, userUUID); err == nil || err != expectedError {
		t.Errorf("Error occurred in PickupLine delete service")
		return
	}
}

func TestGetPickupLineByIdAndUserIdSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	expectedPickupLine := &entities.PickupLine{
		ID: pickupLineUUID,
		User: &entities.User{
			ID: userUUID,
		},
	}

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(expectedPickupLine, nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if retrievedPickupLine, err := service.GetByIdAndUserId(pickupLineUUID, userUUID); err != nil || retrievedPickupLine.ID != pickupLineUUID {
		t.Errorf("Error occurred in PickupLine delete service")
		return
	}
}

func TestGetPickupLineByIdAndUserIdFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	expectedError := errors.New("could not retrieve PickupLine")

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetByIdAndUserId(pickupLineUUID, userUUID); err == nil || err != expectedError {
		t.Errorf("Error occurred in PickupLine delete service")
		return
	}
}

func TestGetPickupLineByIdSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	expectedPickupLine := &entities.PickupLine{
		ID: pickupLineUUID,
	}
	expectedRetrievedPickupLine := &entities.PickupLine{
		ID: pickupLineUUID,
	}

	mockRepository.EXPECT().GetById(pickupLineUUID).Return(expectedPickupLine, nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if retrievedPickupLine, err := service.GetById(pickupLineUUID); err != nil || !retrievedPickupLine.IsEqualTo(expectedRetrievedPickupLine) {
		t.Errorf("Error occurred in PickupLine delete service")
		return
	}
}

func TestGetPickupLineByIdFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	expectedError := errors.New("could not retrive PickupLine")

	mockRepository.EXPECT().GetById(pickupLineUUID).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetById(pickupLineUUID); err == nil || err != expectedError {
		t.Errorf("Error occurred in PickupLine delete service")
		return
	}
}

func TestGetListSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()
	pickupLineList := []*responses.SinglePickupLineInfoResponse{}

	filters := &filters.PickupLineQueryFilters{}
	response := &responses.ElasticSearchPickupLineResponse{
		Total: len(pickupLineList),
		Users: pickupLineList,
	}
	mockElasticWrapper.EXPECT().SearchPickupLines(userUUID, filters).Return(response, nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetList(userUUID, filters); err != nil {
		t.Errorf("Error occurred in PickupLine GetList service")
		return
	}
}

func TestGetListFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()

	filters := &filters.PickupLineQueryFilters{}
	expectedError := errors.New("could not retrieve PickupLine List")
	mockElasticWrapper.EXPECT().SearchPickupLines(userUUID, filters).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetList(userUUID, filters); err == nil || expectedError != err {
		t.Errorf("Error occurred in PickupLine GetList service")
		return
	}
}

func TestGetFeedSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()
	pickupLineList := []*responses.SinglePickupLineInfoResponse{}

	filters := &filters.PickupLineQueryFilters{}
	response := &responses.ElasticSearchPickupLineResponse{
		Total: len(pickupLineList),
		Users: pickupLineList,
	}
	mockElasticWrapper.EXPECT().GetPickupLineFeed(userUUID, filters).Return(response, nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetFeed(userUUID, filters); err != nil {
		t.Errorf("Error occurred in PickupLine Feed service")
		return
	}
}

func TestGetFeedFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()

	filters := &filters.PickupLineQueryFilters{}
	expectedError := errors.New("could not retrieve PickupLine Feed")
	mockElasticWrapper.EXPECT().GetPickupLineFeed(userUUID, filters).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetFeed(userUUID, filters); err == nil || expectedError != err {
		t.Errorf("Error occurred in PickupLine Feed service")
		return
	}
}

func TestUpdateReactionByUserSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()
	pickupLineUUID, _ := utils.GetRandomUUID()
	oldReaction := &entities.Reaction{
		PickupLineId: pickupLineUUID,
		UserId:       userUUID,
		Starred:      false,
		Vote:         entities.Downvote,
	}

	newReaction := &entities.Reaction{
		PickupLineId: pickupLineUUID,
		UserId:       userUUID,
		Starred:      false,
		Vote:         entities.Upvote,
	}

	mockRepository.EXPECT().GetReactionByPickupLineIdAndUserId(pickupLineUUID, userUUID).Return(oldReaction, nil)
	mockRepository.EXPECT().UpdateReaction(newReaction).Return(nil)
	mockElasticWrapper.EXPECT().UpdateUserReaction(userUUID, pickupLineUUID, newReaction, oldReaction).Return(nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if err := service.UpdateReactionByUser(userUUID, pickupLineUUID, newReaction); err != nil {
		t.Errorf("Error occurred in PickupLine Reaction service")
		return
	}
}

func TestUpdateReactionByUserUpdateElasticFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()
	pickupLineUUID, _ := utils.GetRandomUUID()
	oldReaction := &entities.Reaction{
		PickupLineId: pickupLineUUID,
		UserId:       userUUID,
		Starred:      false,
		Vote:         entities.Downvote,
	}

	newReaction := &entities.Reaction{
		PickupLineId: pickupLineUUID,
		UserId:       userUUID,
		Starred:      false,
		Vote:         entities.Upvote,
	}
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetReactionByPickupLineIdAndUserId(pickupLineUUID, userUUID).Return(oldReaction, nil)
	mockRepository.EXPECT().UpdateReaction(newReaction).Return(nil)
	mockElasticWrapper.EXPECT().UpdateUserReaction(userUUID, pickupLineUUID, newReaction, oldReaction).Return(expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if err := service.UpdateReactionByUser(userUUID, pickupLineUUID, newReaction); err == nil || err != expectedError {
		t.Errorf("Error occurred in PickupLine Reaction service")
		return
	}
}

func TestUpdateReactionByUserUpdateReactionFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()
	pickupLineUUID, _ := utils.GetRandomUUID()
	oldReaction := &entities.Reaction{
		PickupLineId: pickupLineUUID,
		UserId:       userUUID,
		Starred:      false,
		Vote:         entities.Downvote,
	}

	newReaction := &entities.Reaction{
		PickupLineId: pickupLineUUID,
		UserId:       userUUID,
		Starred:      false,
		Vote:         entities.Upvote,
	}
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetReactionByPickupLineIdAndUserId(pickupLineUUID, userUUID).Return(oldReaction, nil)
	mockRepository.EXPECT().UpdateReaction(newReaction).Return(expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if err := service.UpdateReactionByUser(userUUID, pickupLineUUID, newReaction); err == nil || err != expectedError {
		t.Errorf("Error occurred in PickupLine Reaction service")
		return
	}
}

func TestUpdateReactionByUserGetOldReactionFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()
	pickupLineUUID, _ := utils.GetRandomUUID()
	newReaction := &entities.Reaction{
		PickupLineId: pickupLineUUID,
		UserId:       userUUID,
		Starred:      false,
		Vote:         entities.Upvote,
	}
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetReactionByPickupLineIdAndUserId(pickupLineUUID, userUUID).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if err := service.UpdateReactionByUser(userUUID, pickupLineUUID, newReaction); err == nil || err != expectedError {
		t.Errorf("Error occurred in PickupLine Reaction service")
		return
	}
}

func TestGetReactionByIdAndUserIdSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()
	pickupLineUUID, _ := utils.GetRandomUUID()
	expectedReaction := &entities.Reaction{
		PickupLineId: pickupLineUUID,
		UserId:       userUUID,
		Starred:      false,
		Vote:         entities.Downvote,
	}

	mockRepository.EXPECT().GetReactionByPickupLineIdAndUserId(pickupLineUUID, userUUID).Return(expectedReaction, nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetReactionByIdAndUserId(pickupLineUUID, userUUID); err != nil {
		t.Errorf("Error occurred in PickupLine Reaction service")
		return
	}
}

func TestGetReactionByIdAndUserIdFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()
	pickupLineUUID, _ := utils.GetRandomUUID()
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetReactionByPickupLineIdAndUserId(pickupLineUUID, userUUID).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetReactionByIdAndUserId(pickupLineUUID, userUUID); err == nil || err != expectedError {
		t.Errorf("Error occurred in PickupLine Reaction service")
		return
	}
}

func TestGetStatisticByIdAndUserIdSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	expectedStatistic := &entities.Statistic{
		NumberOfSuccesses: 0,
		NumberOfFailures:  0,
		NumberOfTries:     0,
		SuccessPercentage: 0,
	}

	mockRepository.EXPECT().GetStatisticsById(pickupLineUUID).Return(expectedStatistic, nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetStatisticsById(pickupLineUUID); err != nil {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestGetStatisticByIdAndUserIdFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetStatisticsById(pickupLineUUID).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetStatisticsById(pickupLineUUID); err == nil || err != expectedError {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestDeleteByUserSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()

	mockElasticWrapper.EXPECT().DeleteUserPickupLines(userUUID).Return(nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if err := service.DeleteByUser(userUUID); err != nil {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestDeleteByUserFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	userUUID, _ := utils.GetRandomUUID()
	expectedError := errors.New("error")

	mockElasticWrapper.EXPECT().DeleteUserPickupLines(userUUID).Return(expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if err := service.DeleteByUser(userUUID); err == nil || err != expectedError {
		t.Errorf("Error occurred in PickupLine pickupLine service")
		return
	}
}

func TestCanUserSeePickupLineTrueVisiblePickupLine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	requestingUserUUID, _ := utils.GetRandomUUID()
	expectedPickupLine := &entities.PickupLine{
		Visible: true,
		User: &entities.User{
			ID: userUUID,
		},
	}

	mockRepository.EXPECT().GetPickupLineForVisibilityCheck(pickupLineUUID).Return(expectedPickupLine, nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if ok := service.CanUserSeePickupLine(pickupLineUUID, requestingUserUUID); !ok {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestCanUserSeePickupLineTrueVisiblePickupLineOfSameUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	expectedPickupLine := &entities.PickupLine{
		Visible: true,
		User: &entities.User{
			ID: userUUID,
		},
	}

	mockRepository.EXPECT().GetPickupLineForVisibilityCheck(pickupLineUUID).Return(expectedPickupLine, nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if ok := service.CanUserSeePickupLine(pickupLineUUID, userUUID); !ok {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestCanUserSeePickupLineTrueNotVisiblePickupLineOfSameUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	expectedPickupLine := &entities.PickupLine{
		Visible: false,
		User: &entities.User{
			ID: userUUID,
		},
	}

	mockRepository.EXPECT().GetPickupLineForVisibilityCheck(pickupLineUUID).Return(expectedPickupLine, nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if ok := service.CanUserSeePickupLine(pickupLineUUID, userUUID); !ok {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestCanUserSeePickupLineFalseNotVisiblePickupLine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	requestinUserUUID, _ := utils.GetRandomUUID()
	expectedPickupLine := &entities.PickupLine{
		Visible: false,
		User: &entities.User{
			ID: userUUID,
		},
	}

	mockRepository.EXPECT().GetPickupLineForVisibilityCheck(pickupLineUUID).Return(expectedPickupLine, nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if ok := service.CanUserSeePickupLine(pickupLineUUID, requestinUserUUID); ok {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestCanUserSeePickupLineFalseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	requestinUserUUID, _ := utils.GetRandomUUID()
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetPickupLineForVisibilityCheck(pickupLineUUID).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if ok := service.CanUserSeePickupLine(pickupLineUUID, requestinUserUUID); ok {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestUpdatePickupLineNoTagsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	randomPickupLineBody := getRandomPickupLineBodyRequestNoTags()
	pickupLineToUpdate := &entities.PickupLine{
		ID: pickupLineUUID,
	}
	expectedUpdatedPickupLine := &entities.PickupLine{
		ID:      pickupLineUUID,
		Title:   randomPickupLineBody.Title,
		Content: randomPickupLineBody.Content,
		Visible: randomPickupLineBody.Visible,
	}

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(pickupLineToUpdate, nil)
	mockRepository.EXPECT().Update(pickupLineToUpdate).Return(nil)
	mockElasticWrapper.EXPECT().UpdatePickupLine(pickupLineToUpdate).Return(nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if updatedPickupLine, err := service.Update(pickupLineUUID, randomPickupLineBody, userUUID); err != nil || !updatedPickupLine.IsEqualTo(expectedUpdatedPickupLine) {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestUpdatePickupLineWithTagsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	randomTagsIds := getRandomTagIdRequests(3)
	tagSearchList := make([][]uuid.UUID, 0, len(randomTagsIds))
	retrievedTagList := make([]*entities.Tag, 0, len(randomTagsIds))
	for _, tagRequest := range randomTagsIds {
		tagUUID, err := uuid.Parse(tagRequest.ID)
		if err != nil {
			t.Errorf("Error occurred in PikcupLine creation service")
			return
		}

		tagSearchList = append(tagSearchList, []uuid.UUID{tagUUID, userUUID})
		retrievedTagList = append(retrievedTagList, &entities.Tag{ID: tagUUID})
	}
	randomPickupLineBodyRequest := getRandomPickupLineBodyRequestNoTags()
	randomPickupLineBodyRequest.Tags = randomTagsIds
	pickupLineToUpdate := &entities.PickupLine{
		ID: pickupLineUUID,
	}
	expectedUpdatedPickupLine := &entities.PickupLine{
		ID:      pickupLineUUID,
		Title:   randomPickupLineBodyRequest.Title,
		Content: randomPickupLineBodyRequest.Content,
		Visible: randomPickupLineBodyRequest.Visible,
		Tags:    retrievedTagList,
	}

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(pickupLineToUpdate, nil)
	mockTagService.EXPECT().GetListByTagIdAndUserIdList(tagSearchList).Return(retrievedTagList, nil)
	mockRepository.EXPECT().Update(pickupLineToUpdate).Return(nil)
	mockElasticWrapper.EXPECT().UpdatePickupLine(pickupLineToUpdate).Return(nil)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if updatedPickupLine, err := service.Update(pickupLineUUID, randomPickupLineBodyRequest, userUUID); err != nil || !updatedPickupLine.IsEqualTo(expectedUpdatedPickupLine) {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestUpdatePickupLineWithTagsFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	randomTagsIds := getRandomTagIdRequests(3)
	tagSearchList := make([][]uuid.UUID, 0, len(randomTagsIds))
	retrievedTagList := make([]*entities.Tag, 0, len(randomTagsIds))
	for _, tagRequest := range randomTagsIds {
		tagUUID, err := uuid.Parse(tagRequest.ID)
		if err != nil {
			t.Errorf("Error occurred in PikcupLine creation service")
			return
		}

		tagSearchList = append(tagSearchList, []uuid.UUID{tagUUID, userUUID})
		retrievedTagList = append(retrievedTagList, &entities.Tag{ID: tagUUID})
	}
	randomPickupLineBodyRequest := getRandomPickupLineBodyRequestNoTags()
	randomPickupLineBodyRequest.Tags = randomTagsIds
	pickupLineToUpdate := &entities.PickupLine{
		ID: pickupLineUUID,
	}
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(pickupLineToUpdate, nil)
	mockTagService.EXPECT().GetListByTagIdAndUserIdList(tagSearchList).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if updatedPickupLine, err := service.Update(pickupLineUUID, randomPickupLineBodyRequest, userUUID); err == nil || updatedPickupLine != nil {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestUpdatePickupLineNoTagsElasticFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	randomPickupLineBody := getRandomPickupLineBodyRequestNoTags()
	pickupLineToUpdate := &entities.PickupLine{
		ID: pickupLineUUID,
	}
	expectedUpdatedPickupLine := &entities.PickupLine{
		ID:      pickupLineUUID,
		Title:   randomPickupLineBody.Title,
		Content: randomPickupLineBody.Content,
		Visible: randomPickupLineBody.Visible,
	}
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(pickupLineToUpdate, nil)
	mockRepository.EXPECT().Update(pickupLineToUpdate).Return(nil)
	mockElasticWrapper.EXPECT().UpdatePickupLine(pickupLineToUpdate).Return(expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if updatedPickupLine, err := service.Update(pickupLineUUID, randomPickupLineBody, userUUID); err == nil || !updatedPickupLine.IsEqualTo(expectedUpdatedPickupLine) {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestUpdatePickupLineNoTagsUpdateFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	randomPickupLineBody := getRandomPickupLineBodyRequestNoTags()
	pickupLineToUpdate := &entities.PickupLine{
		ID: pickupLineUUID,
	}
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(pickupLineToUpdate, nil)
	mockRepository.EXPECT().Update(pickupLineToUpdate).Return(expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if updatedPickupLine, err := service.Update(pickupLineUUID, randomPickupLineBody, userUUID); err == nil || updatedPickupLine != nil {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}

func TestUpdatePickupLineNoTagsRetrieveFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)
	mockTagService := mocks.NewMockTagServiceInterface(ctrl)

	pickupLineUUID, _ := utils.GetRandomUUID()
	userUUID, _ := utils.GetRandomUUID()
	randomPickupLineBody := getRandomPickupLineBodyRequestNoTags()
	expectedError := errors.New("error")

	mockRepository.EXPECT().GetByIdAndUserId(pickupLineUUID, userUUID).Return(nil, expectedError)

	service := CreatePickupLineService(mockRepository, mockTagService, mockElasticWrapper, mockDtoMapper)
	if updatedPickupLine, err := service.Update(pickupLineUUID, randomPickupLineBody, userUUID); err == nil || updatedPickupLine != nil {
		t.Errorf("Error occurred in PickupLine Statistic service")
		return
	}
}
