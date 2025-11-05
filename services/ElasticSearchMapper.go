package services

import (
	"encoding/json"
	"errors"
	"playbook/constants"
	"playbook/elastic_dtos"
	"playbook/entities"
	"playbook/repositories"
	"playbook/responses"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ElasticSearchMapper struct {
	tagRepository        repositories.TagRepositoryInterface
	pickupLineRepository repositories.PickupLineRepositoryInterface
}

type ElasticSearchMapperInterface interface {
	elasticDtoToTags(tagIds []uuid.UUID, userUUID uuid.UUID) ([]*entities.Tag, error)
	hydratePickupLineListFromElasticSearchResponse(response *search.Response, userId uuid.UUID) (*responses.ElasticSearchPickupLineResponse, error)
	getReaction(pickupLineId uuid.UUID, userId uuid.UUID) (*entities.Reaction, error)
	pickupLineToElasticDTO(pickupLine *entities.PickupLine) *elastic_dtos.PickupLineElasticDTO
	pickupLineToUpdateElasticDTO(pickupLine *entities.PickupLine) *elastic_dtos.UpdatePickupLineElasticDTO
	tagToElasticDTO(tag *entities.Tag) *elastic_dtos.TagElasticDTO
	userToElasticDTO(user *entities.User) *elastic_dtos.UserElasticDTO
	hydrateUserListFromElasticSearchResponse(response *search.Response) (*responses.ElasticSearchUserListResponse, error)
}

func GetNewElasticSearchMapper(tagRepository repositories.TagRepositoryInterface, pickupLineRepository repositories.PickupLineRepositoryInterface) ElasticSearchMapperInterface {
	return &ElasticSearchMapper{tagRepository: tagRepository, pickupLineRepository: pickupLineRepository}
}

func (e *ElasticSearchMapper) elasticDtoToTags(tagIds []uuid.UUID, userUUID uuid.UUID) ([]*entities.Tag, error) {
	tagsLen := len(tagIds)
	tagsToSearch := make([][]uuid.UUID, 0, tagsLen)
	tagsCache := make(map[string]*entities.Tag, tagsLen)
	tags := make([]*entities.Tag, 0, tagsLen)
	for _, tagId := range tagIds {
		tag, ok := tagsCache[tagId.String()]
		if !ok {
			tagsToSearch = append(tagsToSearch, []uuid.UUID{tagId, userUUID})
			continue
		}
		tags = append(tags, tag)
	}

	if len(tagsToSearch) != 0 {
		list, err := e.tagRepository.GetListByTagIdAndUserIdList(tagsToSearch)
		if err != nil {
			return nil, err
		}

		for _, tag := range list {
			tagsCache[tag.ID.String()] = tag
			tags = append(tags, tag)
		}
	}

	return tags, nil
}

func (e *ElasticSearchMapper) hydratePickupLineListFromElasticSearchResponse(response *search.Response, userId uuid.UUID) (*responses.ElasticSearchPickupLineResponse, error) {
	pickupLineList := make([]*responses.SinglePickupLineInfoResponse, 0, constants.ITEMS_FOR_PAGE)
	for _, hit := range response.Hits.Hits {
		var pickupLineDTO elastic_dtos.PickupLineElasticDTO
		err := json.Unmarshal(hit.Source_, &pickupLineDTO)
		if err != nil {
			return nil, err
		}

		tags, err := e.elasticDtoToTags(pickupLineDTO.Tags, pickupLineDTO.UserId)
		if err != nil {
			return nil, err
		}
		reaction, err := e.getReaction(pickupLineDTO.Id, userId)
		if err != nil {
			return nil, err
		}

		pickupLineItem := elasticDtoToPickupLine(&pickupLineDTO, tags, reaction)
		pickupLineResponseItem := responses.NewSinglePickupLineInfoResponse(pickupLineItem)
		pickupLineList = append(pickupLineList, pickupLineResponseItem)
	}

	return &responses.ElasticSearchPickupLineResponse{
		Total: int(response.Hits.Total.Value),
		Users: pickupLineList,
	}, nil
}

func (e *ElasticSearchMapper) hydrateUserListFromElasticSearchResponse(response *search.Response) (*responses.ElasticSearchUserListResponse, error) {
	userList := make([]*responses.BaseUserInfoResponse, 0, constants.ITEMS_FOR_PAGE)
	for _, hit := range response.Hits.Hits {
		var userDTO elastic_dtos.UserElasticDTO
		err := json.Unmarshal(hit.Source_, &userDTO)
		if err != nil {
			return nil, err
		}

		userItem := elasticDtoToUser(&userDTO)
		userResponseItem := responses.NewBaseUserInfoResponse(userItem)
		userList = append(userList, userResponseItem)
	}

	return &responses.ElasticSearchUserListResponse{
		Total: int(response.Hits.Total.Value),
		Users: userList,
	}, nil
}

func (e *ElasticSearchMapper) getReaction(pickupLineId uuid.UUID, userId uuid.UUID) (*entities.Reaction, error) {
	reaction, err := e.pickupLineRepository.GetReactionByPickupLineIdAndUserId(pickupLineId, userId)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		reaction.Vote = entities.None
	}
	return reaction, nil
}

func (e *ElasticSearchMapper) pickupLineToElasticDTO(pickupLine *entities.PickupLine) *elastic_dtos.PickupLineElasticDTO {
	tags := make([]uuid.UUID, 0, len(pickupLine.Tags))
	for _, tag := range pickupLine.Tags {
		tags = append(tags, tag.ID)
	}

	return &elastic_dtos.PickupLineElasticDTO{
		Id:          pickupLine.ID,
		Title:       pickupLine.Title,
		Content:     pickupLine.Content,
		Tags:        tags,
		UserId:      pickupLine.User.ID,
		Username:    pickupLine.User.Username,
		DisplayName: pickupLine.User.DisplayName,
		Visible:     pickupLine.Visible,
		UpdatedAt:   pickupLine.UpdatedAt,
	}
}

func (e *ElasticSearchMapper) pickupLineToUpdateElasticDTO(pickupLine *entities.PickupLine) *elastic_dtos.UpdatePickupLineElasticDTO {
	tags := make([]uuid.UUID, 0, len(pickupLine.Tags))
	for _, tag := range pickupLine.Tags {
		tags = append(tags, tag.ID)
	}

	return &elastic_dtos.UpdatePickupLineElasticDTO{
		Title:     pickupLine.Title,
		Content:   pickupLine.Content,
		Tags:      tags,
		Visible:   pickupLine.Visible,
		UpdatedAt: pickupLine.UpdatedAt,
	}
}

func elasticDtoToPickupLine(elasticDTO *elastic_dtos.PickupLineElasticDTO, tags []*entities.Tag, reaction *entities.Reaction) *entities.PickupLine {
	return &entities.PickupLine{
		ID:      elasticDTO.Id,
		Title:   elasticDTO.Title,
		Content: elasticDTO.Content,
		Visible: elasticDTO.Visible,
		User: &entities.User{
			ID:          elasticDTO.UserId,
			Username:    elasticDTO.Username,
			DisplayName: elasticDTO.DisplayName,
		},
		UserReaction: reaction,
		Tags:         tags,
		Statistics: &entities.Statistic{
			NumberOfSuccesses: elasticDTO.NumberOfSuccesses,
			NumberOfTries:     elasticDTO.NumberOfTries,
			NumberOfFailures:  elasticDTO.NumberOfFailures,
			SuccessPercentage: elasticDTO.SuccessPercentage,
		},
		UpdatedAt: elasticDTO.UpdatedAt,
	}
}

func (e *ElasticSearchMapper) tagToElasticDTO(tag *entities.Tag) *elastic_dtos.TagElasticDTO {
	return &elastic_dtos.TagElasticDTO{
		Id:     tag.ID,
		Name:   tag.Name,
		UserId: tag.UserId,
	}
}

func (e *ElasticSearchMapper) userToElasticDTO(user *entities.User) *elastic_dtos.UserElasticDTO {
	return &elastic_dtos.UserElasticDTO{
		Id:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
	}
}

func elasticDtoToUser(elasticDTO *elastic_dtos.UserElasticDTO) *entities.User {
	return &entities.User{
		ID:          elasticDTO.Id,
		DisplayName: elasticDTO.DisplayName,
		Username:    elasticDTO.Username,
	}
}
