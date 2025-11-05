package services

import (
	"errors"
	"log"
	"playbook/entities"
	"playbook/filters"
	"playbook/mappers"
	"playbook/repositories"
	"playbook/requests"
	"playbook/responses"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PickupLineService struct {
	pickupLineRepository repositories.PickupLineRepositoryInterface
	tagService           TagServiceInterface
	elasticSearchWrapper ElasticSearchWrapperInterface
	dtoMapper            mappers.DtoMapperInterface
}

func (u *PickupLineService) Create(pickupLineDTO *requests.PickupLineBodyRequest, user *entities.User) (*entities.PickupLine, error) {
	pickupLineEntity := u.dtoMapper.GetPickupLineEntityFromDTO(pickupLineDTO)
	pickupLineEntity.User = user
	tags, err := u.GetTagsFromRequest(pickupLineDTO.Tags, user.ID)
	if err != nil {
		return nil, err
	}
	pickupLineEntity.Tags = tags
	pickupLineEntity.Statistics = &entities.Statistic{}

	err = u.pickupLineRepository.Create(pickupLineEntity)
	if err != nil {
		return nil, err
	}

	err = u.elasticSearchWrapper.IndexPickupLine(pickupLineEntity)
	if err != nil { //todo: logging
		u.elasticSearchWrapper.DeletePickupLine(pickupLineEntity.ID)
		u.pickupLineRepository.Delete(pickupLineEntity)
		return nil, err
	}

	return pickupLineEntity, nil
}

func (u *PickupLineService) GetTagsFromRequest(rawTags []*requests.TagIdRequest, userUUID uuid.UUID) ([]*entities.Tag, error) {
	if len(rawTags) == 0 {
		return nil, nil
	}
	tagSearchList := make([][]uuid.UUID, 0, len(rawTags))
	for _, tag := range rawTags {
		tagUUID, err := uuid.Parse(tag.ID)
		if err != nil {
			return nil, err
		}

		tagTuple := []uuid.UUID{tagUUID, userUUID}
		tagSearchList = append(tagSearchList, tagTuple)
	}

	tags, err := u.tagService.GetListByTagIdAndUserIdList(tagSearchList)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (u *PickupLineService) Delete(pickupLineUUID uuid.UUID, userUUID uuid.UUID) error {
	pickupLineToDelete, err := u.pickupLineRepository.GetByIdAndUserId(pickupLineUUID, userUUID)
	if err != nil {
		return err
	}
	err = u.elasticSearchWrapper.DeletePickupLine(pickupLineToDelete.ID)
	if err != nil {
		return err
	}
	return u.pickupLineRepository.Delete(pickupLineToDelete)
}

func (u *PickupLineService) GetByIdAndUserId(pickupLineId uuid.UUID, userUUID uuid.UUID) (*entities.PickupLine, error) {
	pickupLineEntity, err := u.pickupLineRepository.GetByIdAndUserId(pickupLineId, userUUID)
	if err != nil {
		return nil, err
	}
	return pickupLineEntity, nil
}

func (u *PickupLineService) GetById(pickupLineId uuid.UUID) (*entities.PickupLine, error) {
	return u.pickupLineRepository.GetById(pickupLineId)
}

func (u *PickupLineService) GetList(userId uuid.UUID, filters *filters.PickupLineQueryFilters) (*responses.ElasticSearchPickupLineResponse, error) {
	return u.elasticSearchWrapper.SearchPickupLines(userId, filters)
}

func (u *PickupLineService) GetFeed(userId uuid.UUID, filters *filters.PickupLineQueryFilters) (*responses.ElasticSearchPickupLineResponse, error) {
	return u.elasticSearchWrapper.GetPickupLineFeed(userId, filters)
}

func (u *PickupLineService) UpdateReactionByUser(userUUID uuid.UUID, pickupLineUUID uuid.UUID, newReaction *entities.Reaction) error {
	oldReaction, err := u.pickupLineRepository.GetReactionByPickupLineIdAndUserId(pickupLineUUID, userUUID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	err = u.pickupLineRepository.UpdateReaction(newReaction)
	if err != nil {
		return err
	}

	return u.elasticSearchWrapper.UpdateUserReaction(userUUID, pickupLineUUID, newReaction, oldReaction)
}

func (u *PickupLineService) GetReactionByIdAndUserId(pickupLineUUID uuid.UUID, userUUID uuid.UUID) (*entities.Reaction, error) {
	return u.pickupLineRepository.GetReactionByPickupLineIdAndUserId(pickupLineUUID, userUUID)
}

func (u *PickupLineService) GetStatisticsById(pickupLineId uuid.UUID) (*entities.Statistic, error) {
	return u.pickupLineRepository.GetStatisticsById(pickupLineId)
}

func (u *PickupLineService) Update(pickupLineUUID uuid.UUID, pickupLineDTO *requests.PickupLineBodyRequest, userUUID uuid.UUID) (*entities.PickupLine, error) {
	pickupLineToUpdate, err := u.pickupLineRepository.GetByIdAndUserId(pickupLineUUID, userUUID)
	if err != nil {
		return nil, err
	}

	tags, err := u.GetTagsFromRequest(pickupLineDTO.Tags, userUUID)
	if err != nil {
		return nil, err
	}

	pickupLineToUpdate.Title = pickupLineDTO.Title
	pickupLineToUpdate.Content = pickupLineDTO.Content
	pickupLineToUpdate.Visible = pickupLineDTO.Visible
	pickupLineToUpdate.Tags = tags
	err = u.pickupLineRepository.Update(pickupLineToUpdate)
	if err != nil {
		return nil, err
	}

	return pickupLineToUpdate, u.elasticSearchWrapper.UpdatePickupLine(pickupLineToUpdate)
}

func (u *PickupLineService) DeleteByUser(userUUID uuid.UUID) error {
	return u.elasticSearchWrapper.DeleteUserPickupLines(userUUID)
}

func (u *PickupLineService) CanUserSeePickupLine(pickupLineUUID uuid.UUID, userUUID uuid.UUID) bool {
	pickupLine, err := u.pickupLineRepository.GetPickupLineForVisibilityCheck(pickupLineUUID)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	return pickupLine.Visible || pickupLine.User.ID == userUUID
}

func CreatePickupLineService(pickupLineRepository repositories.PickupLineRepositoryInterface, tagService TagServiceInterface, elasticSearchWrapper ElasticSearchWrapperInterface, dtoMapper mappers.DtoMapperInterface) *PickupLineService {
	return &PickupLineService{pickupLineRepository: pickupLineRepository, tagService: tagService, elasticSearchWrapper: elasticSearchWrapper, dtoMapper: dtoMapper}
}
