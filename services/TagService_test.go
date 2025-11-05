package services

import (
	"errors"
	"playbook/entities"
	"playbook/mocks"
	"playbook/requests"
	"playbook/utils"
	"testing"

	"go.uber.org/mock/gomock"
	"github.com/google/uuid"
)

func getRandomTagBodyRequest() *requests.TagBodyRequest {
	return &requests.TagBodyRequest{
		Name:        utils.GetRandomTagName(),
		Description: utils.GetRandomTagDescription(),
	}
}

func TestAddTagSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	randomTag := getRandomTagBodyRequest()
	tag := &entities.Tag{
		ID:          uuid.Max,
		Description: randomTag.Description,
		Name:        randomTag.Name,
		UserId:      uuid.Max,
	}
	mockDtoMapper.EXPECT().TagBodyRequestToEntity(randomTag, uuid.Max).Return(tag)
	mockElasticWrapper.EXPECT().IndexTag(tag).Return("elasticId", nil)
	mockRepository.EXPECT().Create(tag).Return(nil)

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if createdTag, err := service.Create(randomTag, uuid.Max); err != nil || createdTag.ElasticSearchId != "elasticId" {
		t.Errorf("Error occurred in tag creation service")
		return
	}
}

func TestAddTagElasticSearchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	randomTag := getRandomTagBodyRequest()
	tag := &entities.Tag{
		ID:          uuid.Max,
		Description: randomTag.Description,
		Name:        randomTag.Name,
		UserId:      uuid.Max,
	}
	mockDtoMapper.EXPECT().TagBodyRequestToEntity(randomTag, uuid.Max).Return(tag)
	err := errors.New("elasticSearch index error")
	mockElasticWrapper.EXPECT().IndexTag(tag).Return("", err)

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if _, err := service.Create(randomTag, uuid.Max); err == nil {
		t.Errorf("Error occurred in tag creation service")
		return
	}
}

func TestAddTagDatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	randomTag := getRandomTagBodyRequest()
	tag := &entities.Tag{
		ID:          uuid.Max,
		Description: randomTag.Description,
		Name:        randomTag.Name,
		UserId:      uuid.Max,
	}
	mockDtoMapper.EXPECT().TagBodyRequestToEntity(randomTag, uuid.Max).Return(tag)
	mockElasticWrapper.EXPECT().IndexTag(tag).Return("elasticId", nil)
	err := errors.New("db error")
	mockRepository.EXPECT().Create(tag).Return(err)
	mockElasticWrapper.EXPECT().DeleteTag(tag).Return(nil)

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if _, err := service.Create(randomTag, uuid.Max); err == nil {
		t.Errorf("Error occurred in tag creation service")
		return
	}
}

func TestDeleteTagSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	tag := &entities.Tag{
		ID:     uuid.Max,
		UserId: uuid.Max,
	}
	mockElasticWrapper.EXPECT().DeleteTag(tag).Return(nil)
	mockRepository.EXPECT().Delete(tag).Return(nil)

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if err := service.Delete(tag); err != nil {
		t.Errorf("Error occurred in tag deletion service")
		return
	}
}

func TestDeleteTagElasticSearchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	tag := &entities.Tag{
		ID:     uuid.Max,
		UserId: uuid.Max,
	}
	mockElasticWrapper.EXPECT().DeleteTag(tag).Return(errors.New("elastic delete error"))
	mockRepository.EXPECT().Delete(tag).Return(nil)

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if err := service.Delete(tag); err != nil {
		t.Errorf("Error occurred in tag deletion service")
		return
	}
}

func TestDeleteTagDatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	tag := &entities.Tag{
		ID:     uuid.Max,
		UserId: uuid.Max,
	}
	mockElasticWrapper.EXPECT().DeleteTag(tag).Return(nil)
	mockRepository.EXPECT().Delete(tag).Return(errors.New("elastic delete error"))

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if err := service.Delete(tag); err == nil {
		t.Errorf("Error occurred in tag deletion service")
		return
	}
}

func TestUpdateTagSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	tagRequest := getRandomTagBodyRequest()
	tag := &entities.Tag{
		ID:     uuid.Max,
		UserId: uuid.Max,
	}
	mockRepository.EXPECT().GetById(tag.ID.String()).Return(tag, nil)
	mockElasticWrapper.EXPECT().UpdateTag(tag).Return(nil)
	mockRepository.EXPECT().Update(tag).Return(nil)

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if updatedTag, err := service.Update(tag.ID.String(), tagRequest); err != nil || updatedTag.Description != tagRequest.Description || updatedTag.Name != tagRequest.Name {
		t.Errorf("Error occurred in tag update service")
		return
	}
}

func TestUpdateTagNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	tagRequest := getRandomTagBodyRequest()
	tag := &entities.Tag{
		ID:     uuid.Max,
		UserId: uuid.Max,
	}
	mockRepository.EXPECT().GetById(tag.ID.String()).Return(nil, errors.New("tag not found"))

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if _, err := service.Update(tag.ID.String(), tagRequest); err == nil {
		t.Errorf("Error occurred in tag update service")
		return
	}
}

func TestUpdateTagDatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	tagRequest := getRandomTagBodyRequest()
	tag := &entities.Tag{
		ID:     uuid.Max,
		UserId: uuid.Max,
	}
	mockRepository.EXPECT().GetById(tag.ID.String()).Return(tag, nil)
	mockRepository.EXPECT().Update(tag).Return(errors.New("database error"))

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if _, err := service.Update(tag.ID.String(), tagRequest); err == nil {
		t.Errorf("Error occurred in tag update service")
		return
	}
}

func TestUpdateTagElasticSearchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	tagRequest := getRandomTagBodyRequest()
	tag := &entities.Tag{
		ID:     uuid.Max,
		UserId: uuid.Max,
	}
	mockRepository.EXPECT().GetById(tag.ID.String()).Return(tag, nil)
	mockRepository.EXPECT().Update(tag).Return(nil)
	mockElasticWrapper.EXPECT().UpdateTag(tag).Return(errors.New("elastic error"))

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if _, err := service.Update(tag.ID.String(), tagRequest); err == nil {
		t.Errorf("Error occurred in tag update service")
		return
	}
}

func TestGetTagListByUserIdSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	tagList := make([]*entities.Tag, 0, 0)
	mockRepository.EXPECT().GetByUserId(uuid.Max).Return(tagList, nil)

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetListByUserId(uuid.Max); err != nil {
		t.Errorf("Error occurred in get tag list service")
		return
	}
}

func TestGetTagListByUserIdError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	mockRepository.EXPECT().GetByUserId(uuid.Max).Return(nil, errors.New("database error"))

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetListByUserId(uuid.Max); err == nil {
		t.Errorf("Error occurred in get tag list service")
		return
	}
}

func TestGetTagListByTagIdAndByUserIdSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	tagList := make([]*entities.Tag, 0, 0)
	uuidsTuples := make([][]uuid.UUID, 0, 0)
	mockRepository.EXPECT().GetListByTagIdAndUserIdList(uuidsTuples).Return(tagList, nil)

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetListByTagIdAndUserIdList(uuidsTuples); err != nil {
		t.Errorf("Error occurred in get tag list by tag and user service")
		return
	}
}

func TestGetTagListByTagIdAndByUserIdError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	mockElasticWrapper := mocks.NewMockElasticSearchWrapperInterface(ctrl)
	mockDtoMapper := mocks.NewMockDtoMapperInterface(ctrl)

	uuidsTuples := make([][]uuid.UUID, 0, 0)
	mockRepository.EXPECT().GetListByTagIdAndUserIdList(uuidsTuples).Return(nil, errors.New("database error"))

	service := CreateTagService(mockRepository, mockElasticWrapper, mockDtoMapper)
	if _, err := service.GetListByTagIdAndUserIdList(uuidsTuples); err == nil {
		t.Errorf("Error occurred in get tag list by tag and user service")
		return
	}
}
