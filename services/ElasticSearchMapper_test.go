package services

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"playbook/elastic_dtos"
	"playbook/entities"
	"playbook/mocks"
	"playbook/responses"
	"playbook/utils"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"gorm.io/gorm"

	"github.com/google/uuid"
	"go.uber.org/mock/gomock"
)

func TestElasticDtoToTagsEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	tagIds := []uuid.UUID{}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if mappedTags, err := mapper.elasticDtoToTags(tagIds, uuid.Nil); err != nil || mappedTags == nil {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestElasticDtoToTagsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	tagIds := make([]uuid.UUID, 0, 0)

	randomTagsIds := getRandomTagIdRequests(3)
	randomTagsIds = append(randomTagsIds, randomTagsIds...)
	userUUID, _ := utils.GetRandomUUID()
	tagSearchList := make([][]uuid.UUID, 0, len(randomTagsIds))
	retrievedTagList := make([]*entities.Tag, 0, len(randomTagsIds))
	for _, tagRequest := range randomTagsIds {
		tagUUID, err := uuid.Parse(tagRequest.ID)
		if err != nil {
			t.Errorf("Error occurred in PikcupLine creation service")
			return
		}

		tagIds = append(tagIds, tagUUID)
		tagSearchList = append(tagSearchList, []uuid.UUID{tagUUID, userUUID})
		retrievedTagList = append(retrievedTagList, &entities.Tag{ID: tagUUID})
	}

	mockTagRepository.EXPECT().GetListByTagIdAndUserIdList(tagSearchList).Return(retrievedTagList, nil)

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if mappedTags, err := mapper.elasticDtoToTags(tagIds, userUUID); err != nil || mappedTags == nil {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestElasticDtoToTagsFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	tagIds := make([]uuid.UUID, 0, 0)

	randomTagsIds := getRandomTagIdRequests(3)
	randomTagsIds = append(randomTagsIds, randomTagsIds...)
	userUUID, _ := utils.GetRandomUUID()
	tagSearchList := make([][]uuid.UUID, 0, len(randomTagsIds))
	retrievedTagList := make([]*entities.Tag, 0, len(randomTagsIds))
	for _, tagRequest := range randomTagsIds {
		tagUUID, err := uuid.Parse(tagRequest.ID)
		if err != nil {
			t.Errorf("Error occurred in PikcupLine creation service")
			return
		}

		tagIds = append(tagIds, tagUUID)
		tagSearchList = append(tagSearchList, []uuid.UUID{tagUUID, userUUID})
		retrievedTagList = append(retrievedTagList, &entities.Tag{ID: tagUUID})
	}
	expectedError := errors.New("error")

	mockTagRepository.EXPECT().GetListByTagIdAndUserIdList(tagSearchList).Return(nil, expectedError)

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if mappedTags, err := mapper.elasticDtoToTags(tagIds, userUUID); err == nil || mappedTags != nil {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestHydrateUserListFromElasticEmptyResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	response := &search.Response{
		Hits: types.HitsMetadata{
			Hits: []types.Hit{},
			Total: &types.TotalHits{
				Value: 0,
			},
		},
	}
	expectedResponse := &responses.ElasticSearchUserListResponse{
		Total: 0,
		Users: []*responses.BaseUserInfoResponse{},
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if userResponse, err := mapper.hydrateUserListFromElasticSearchResponse(response); err != nil || userResponse.Total != expectedResponse.Total || len(userResponse.Users) != len(expectedResponse.Users) {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestHydrateUserListFromElasticResponseSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const numberOfUsers = 4

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	hits := make([]types.Hit, 0, numberOfUsers)
	userList := make([]*responses.BaseUserInfoResponse, 0, numberOfUsers)
	for i := 0; i < numberOfUsers; i++ {
		userDto := &elastic_dtos.UserElasticDTO{
			Id:          uuid.Max,
			Username:    utils.GetRandomUsername(),
			DisplayName: utils.GetRandomString(20),
		}
		user := &responses.BaseUserInfoResponse{
			ID:          userDto.Id,
			Username:    userDto.Username,
			DisplayName: userDto.DisplayName,
		}

		userList = append(userList, user)
		hitSource, _ := json.Marshal(userDto)
		hit := types.Hit{
			Source_: hitSource,
		}
		hits = append(hits, hit)
	}

	response := &search.Response{
		Hits: types.HitsMetadata{
			Hits: hits,
			Total: &types.TotalHits{
				Value: numberOfUsers,
			},
		},
	}
	expectedResponse := &responses.ElasticSearchUserListResponse{
		Total: numberOfUsers,
		Users: userList,
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if userResponse, err := mapper.hydrateUserListFromElasticSearchResponse(response); err != nil || userResponse.Total != expectedResponse.Total || len(userResponse.Users) != len(expectedResponse.Users) {
		t.Errorf("Error occurred in user elastic mapping service")
		return
	}
}

func TestHydratePickupLineListFromElasticEmptyResponse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	response := &search.Response{
		Hits: types.HitsMetadata{
			Hits: []types.Hit{},
			Total: &types.TotalHits{
				Value: 0,
			},
		},
	}
	expectedResponse := &responses.ElasticSearchPickupLineResponse{
		Total: 0,
		Users: []*responses.SinglePickupLineInfoResponse{},
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if pickupLineResponse, err := mapper.hydratePickupLineListFromElasticSearchResponse(response, uuid.Nil); err != nil || pickupLineResponse.Total != expectedResponse.Total || len(pickupLineResponse.Users) != len(expectedResponse.Users) {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestHydratePickupLineListFromElasticEmptyResponseNoTagsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	response := &search.Response{
		Hits: types.HitsMetadata{
			Hits: []types.Hit{},
			Total: &types.TotalHits{
				Value: 0,
			},
		},
	}
	expectedResponse := &responses.ElasticSearchPickupLineResponse{
		Total: 0,
		Users: []*responses.SinglePickupLineInfoResponse{},
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if pickupLineResponse, err := mapper.hydratePickupLineListFromElasticSearchResponse(response, uuid.Nil); err != nil || pickupLineResponse.Total != expectedResponse.Total || len(pickupLineResponse.Users) != len(expectedResponse.Users) {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestHydratePickupLineListFromElasticResponseNoTagsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const numberOfPickupLines = 4

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	hits := make([]types.Hit, 0, numberOfPickupLines)
	pickupLineList := make([]*responses.SinglePickupLineInfoResponse, 0, numberOfPickupLines)
	for i := 0; i < numberOfPickupLines; i++ {
		pickupLineDto := &elastic_dtos.PickupLineElasticDTO{
			Id:                uuid.Max,
			Title:             utils.GetRandomString(10),
			Content:           utils.GetRandomTagDescription(),
			Tags:              []uuid.UUID{},
			UserId:            uuid.Nil,
			Username:          utils.GetRandomUsername(),
			DisplayName:       utils.GetRandomString(20),
			Visible:           false,
			Starred:           false,
			NumberOfSuccesses: 0,
			NumberOfFailures:  0,
			NumberOfTries:     0,
			SuccessPercentage: 0,
			UpdatedAt:         time.Time{},
		}
		pickupLine := &responses.SinglePickupLineInfoResponse{
			ID:      pickupLineDto.Id,
			Title:   pickupLineDto.Title,
			Content: pickupLineDto.Content,
			User: &responses.BaseUserInfoResponse{
				ID:          pickupLineDto.UserId,
				Username:    pickupLineDto.Username,
				DisplayName: pickupLineDto.DisplayName,
			},
			Tags:         []*entities.Tag{},
			Statistics:   &entities.Statistic{},
			UserReaction: &entities.Reaction{},
			Visible:      false,
			UpdatedAt:    pickupLineDto.UpdatedAt,
		}

		pickupLineList = append(pickupLineList, pickupLine)
		hitSource, _ := json.Marshal(pickupLineDto)
		hit := types.Hit{
			Source_: hitSource,
		}
		if i == 0 {
			mockPickupLineRepository.EXPECT().
				GetReactionByPickupLineIdAndUserId(pickupLine.ID, pickupLineDto.UserId).
				Return(&entities.Reaction{Vote: entities.None}, gorm.ErrRecordNotFound)
		} else {
			mockPickupLineRepository.EXPECT().
				GetReactionByPickupLineIdAndUserId(pickupLine.ID, pickupLineDto.UserId).
				Return(&entities.Reaction{}, nil)
		}

		hits = append(hits, hit)
	}

	response := &search.Response{
		Hits: types.HitsMetadata{
			Hits: hits,
			Total: &types.TotalHits{
				Value: numberOfPickupLines,
			},
		},
	}
	expectedResponse := &responses.ElasticSearchPickupLineResponse{
		Total: numberOfPickupLines,
		Users: pickupLineList,
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if pickupLineResponse, err := mapper.hydratePickupLineListFromElasticSearchResponse(response, uuid.Nil); err != nil || pickupLineResponse.Total != expectedResponse.Total || len(pickupLineResponse.Users) != len(expectedResponse.Users) {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestHydratePickupLineListFromElasticResponseNoTagsInvalidElasticBodyFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const numberOfPickupLines = 4

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	hits := make([]types.Hit, 0, numberOfPickupLines)
	for i := 0; i < numberOfPickupLines; i++ {
		wrongEntity := &entities.User{}
		hitSource, _ := xml.Marshal(wrongEntity)
		hit := types.Hit{
			Source_: hitSource,
		}

		hits = append(hits, hit)
	}

	response := &search.Response{
		Hits: types.HitsMetadata{
			Hits: hits,
			Total: &types.TotalHits{
				Value: numberOfPickupLines,
			},
		},
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if _, err := mapper.hydratePickupLineListFromElasticSearchResponse(response, uuid.Nil); err == nil {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestHydratePickupLineListFromElasticResponseNoTagsReactionErrorFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	const numberOfPickupLines = 4

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	hits := make([]types.Hit, 0, numberOfPickupLines)
	for i := 0; i < numberOfPickupLines; i++ {
		pickupLineDto := &elastic_dtos.PickupLineElasticDTO{
			Id:                uuid.Max,
			Title:             utils.GetRandomString(10),
			Content:           utils.GetRandomTagDescription(),
			Tags:              []uuid.UUID{},
			UserId:            uuid.Nil,
			Username:          utils.GetRandomUsername(),
			DisplayName:       utils.GetRandomString(20),
			Visible:           false,
			Starred:           false,
			NumberOfSuccesses: 0,
			NumberOfFailures:  0,
			NumberOfTries:     0,
			SuccessPercentage: 0,
			UpdatedAt:         time.Time{},
		}
		hitSource, _ := json.Marshal(pickupLineDto)
		hit := types.Hit{
			Source_: hitSource,
		}

		hits = append(hits, hit)
	}

	expectedError := errors.New("error")
	mockPickupLineRepository.EXPECT().GetReactionByPickupLineIdAndUserId(uuid.Max, uuid.Nil).Return(nil, expectedError)

	response := &search.Response{
		Hits: types.HitsMetadata{
			Hits: hits,
			Total: &types.TotalHits{
				Value: numberOfPickupLines,
			},
		},
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if _, err := mapper.hydratePickupLineListFromElasticSearchResponse(response, uuid.Nil); err == nil || err != expectedError {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestPickupLineToUpdateElasticDtoNoTagsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	pickupLine := &entities.PickupLine{
		Title:     utils.GetRandomString(10),
		Content:   utils.GetRandomTagDescription(),
		Tags:      []*entities.Tag{},
		Visible:   false,
		UpdatedAt: time.Time{},
	}

	expectedUpdateDto := &elastic_dtos.UpdatePickupLineElasticDTO{
		Title:     pickupLine.Title,
		Content:   pickupLine.Content,
		Tags:      []uuid.UUID{},
		Visible:   pickupLine.Visible,
		UpdatedAt: time.Time{},
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if pickupLineUpdateDto := mapper.pickupLineToUpdateElasticDTO(pickupLine); pickupLineUpdateDto == nil || expectedUpdateDto.Title != pickupLineUpdateDto.Title {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestPickupLineToUpdateElasticDtoWithTagsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	tagUUID, _ := utils.GetRandomUUID()
	pickupLine := &entities.PickupLine{
		Title:     utils.GetRandomString(10),
		Content:   utils.GetRandomTagDescription(),
		Tags:      []*entities.Tag{&entities.Tag{ID: tagUUID}},
		Visible:   false,
		UpdatedAt: time.Time{},
	}

	expectedUpdateDto := &elastic_dtos.UpdatePickupLineElasticDTO{
		Title:     pickupLine.Title,
		Content:   pickupLine.Content,
		Tags:      []uuid.UUID{tagUUID},
		Visible:   pickupLine.Visible,
		UpdatedAt: time.Time{},
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if pickupLineUpdateDto := mapper.pickupLineToUpdateElasticDTO(pickupLine); pickupLineUpdateDto == nil || expectedUpdateDto.Title != pickupLineUpdateDto.Title || pickupLineUpdateDto.Tags[0] != tagUUID {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestPickupLineToElasticDtoNoTagsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	pickupLineUUID, _ := utils.GetRandomUUID()
	pickupLine := &entities.PickupLine{
		ID:      pickupLineUUID,
		Title:   utils.GetRandomString(10),
		Content: utils.GetRandomTagDescription(),
		User: &entities.User{
			ID:          uuid.MustParse("bc1fbbe9-2dcd-4857-aa2d-fc152fe01d18"),
			Username:    utils.GetRandomUsername(),
			DisplayName: utils.GetRandomString(10),
		},
		Tags:         []*entities.Tag{},
		Statistics:   &entities.Statistic{},
		Reactions:    []*entities.Reaction{},
		UserReaction: &entities.Reaction{},
		Visible:      false,
		UpdatedAt:    time.Time{},
	}

	expectedUpdateDto := &elastic_dtos.PickupLineElasticDTO{
		Id:          pickupLineUUID,
		Title:       pickupLine.Title,
		Content:     pickupLine.Content,
		Tags:        []uuid.UUID{},
		UserId:      uuid.MustParse("bc1fbbe9-2dcd-4857-aa2d-fc152fe01d18"),
		Username:    pickupLine.User.Username,
		DisplayName: pickupLine.User.DisplayName,
		Visible:     pickupLine.Visible,
		UpdatedAt:   time.Time{},
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if pickupLineUpdateDto := mapper.pickupLineToElasticDTO(pickupLine); pickupLineUpdateDto == nil || expectedUpdateDto.Title != pickupLineUpdateDto.Title {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestPickupLineToElasticDtoWithTagsSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	tagUUID, _ := utils.GetRandomUUID()
	pickupLineUUID, _ := utils.GetRandomUUID()
	pickupLine := &entities.PickupLine{
		ID:      pickupLineUUID,
		Title:   utils.GetRandomString(10),
		Content: utils.GetRandomTagDescription(),
		User: &entities.User{
			ID:          uuid.MustParse("bc1fbbe9-2dcd-4857-aa2d-fc152fe01d18"),
			Username:    utils.GetRandomUsername(),
			DisplayName: utils.GetRandomString(10),
		},
		Tags:         []*entities.Tag{&entities.Tag{ID: tagUUID}},
		Statistics:   &entities.Statistic{},
		Reactions:    []*entities.Reaction{},
		UserReaction: &entities.Reaction{},
		Visible:      utils.GetRandomBool(),
		UpdatedAt:    time.Time{},
	}

	expectedUpdateDto := &elastic_dtos.PickupLineElasticDTO{
		Id:          pickupLineUUID,
		Title:       pickupLine.Title,
		Content:     pickupLine.Content,
		Tags:        []uuid.UUID{tagUUID},
		UserId:      uuid.MustParse("bc1fbbe9-2dcd-4857-aa2d-fc152fe01d18"),
		Username:    pickupLine.User.Username,
		DisplayName: pickupLine.User.DisplayName,
		Visible:     pickupLine.Visible,
		UpdatedAt:   time.Time{},
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if pickupLineUpdateDto := mapper.pickupLineToElasticDTO(pickupLine); pickupLineUpdateDto == nil || expectedUpdateDto.Title != pickupLineUpdateDto.Title || pickupLineUpdateDto.Tags[0] != tagUUID {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestTagToElasticDto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	tagUUID, _ := utils.GetRandomUUID()
	tag := &entities.Tag{
		ID:          tagUUID,
		Name:        utils.GetRandomTagName(),
		Description: utils.GetRandomTagDescription(),
		UserId:      uuid.Max,
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if tagDto := mapper.tagToElasticDTO(tag); tagDto == nil || tagDto.Name != tag.Name || tagDto.Id != tag.ID || tagDto.UserId != tag.UserId {
		t.Errorf("Error occurred in PikcupLine creation service")
		return
	}
}

func TestUserToElasticDto(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPickupLineRepository := mocks.NewMockPickupLineRepositoryInterface(ctrl)
	mockTagRepository := mocks.NewMockTagRepositoryInterface(ctrl)
	userUUID, _ := utils.GetRandomUUID()
	user := &entities.User{
		ID:          userUUID,
		DisplayName: utils.GetRandomString(10),
		Username:    utils.GetRandomUsername(),
	}

	mapper := GetNewElasticSearchMapper(mockTagRepository, mockPickupLineRepository)
	if userDto := mapper.userToElasticDTO(user); userDto == nil || userDto.Id != user.ID || userDto.Username != user.Username || userDto.DisplayName != user.DisplayName {
		t.Errorf("Error occurred in user dto mapping")
		return
	}
}
