package services

import (
	"playbook/entities"
	"playbook/mappers"
	"playbook/repositories"
	"playbook/requests"

	"github.com/google/uuid"
)

type TagServiceInterface interface {
	Create(tag *requests.TagBodyRequest, userUUID uuid.UUID) (*entities.Tag, error)
	Delete(tag *entities.Tag) error
	Update(tagId string, tag *requests.TagBodyRequest) (*entities.Tag, error)
	GetListByUserId(userId uuid.UUID) ([]*entities.Tag, error)
	GetListByTagIdAndUserIdList(tagList [][]uuid.UUID) ([]*entities.Tag, error)
}

type TagService struct {
	tagRepository        repositories.TagRepositoryInterface
	elasticSearchWrapper ElasticSearchWrapperInterface
	dtoMapper            mappers.DtoMapperInterface
}

func (u *TagService) Create(tag *requests.TagBodyRequest, userUUID uuid.UUID) (*entities.Tag, error) {
	tagToCreate := u.dtoMapper.TagBodyRequestToEntity(tag, userUUID)
	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	tagToCreate.ID = uuid

	elasticSearchId, err := u.elasticSearchWrapper.IndexTag(tagToCreate)
	if err != nil {
		return nil, err
	}
	tagToCreate.ElasticSearchId = elasticSearchId

	err = u.tagRepository.Create(tagToCreate)
	if err != nil {
		u.elasticSearchWrapper.DeleteTag(tagToCreate)
		return nil, err
	}

	return tagToCreate, nil
}

func (u *TagService) Delete(tag *entities.Tag) error {
	_ = u.elasticSearchWrapper.DeleteTag(tag)
	return u.tagRepository.Delete(tag)
}

func (u *TagService) Update(tagId string, tag *requests.TagBodyRequest) (*entities.Tag, error) {
	tagToUpdate, err := u.tagRepository.GetById(tagId)
	if err != nil {
		return nil, err
	}

	tagToUpdate.Name = tag.Name
	tagToUpdate.Description = tag.Description
	if err = u.tagRepository.Update(tagToUpdate); err != nil {
		return nil, err
	}

	return tagToUpdate, u.elasticSearchWrapper.UpdateTag(tagToUpdate)
}

func CreateTagService(tagRepository repositories.TagRepositoryInterface, elasticSearchWrapper ElasticSearchWrapperInterface, dtoMapper mappers.DtoMapperInterface) TagServiceInterface {
	return &TagService{tagRepository: tagRepository, elasticSearchWrapper: elasticSearchWrapper, dtoMapper: dtoMapper}
}

func (u *TagService) GetListByUserId(userId uuid.UUID) ([]*entities.Tag, error) {
	tags, err := u.tagRepository.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (u *TagService) GetListByTagIdAndUserIdList(tagList [][]uuid.UUID) ([]*entities.Tag, error) {
	tags, err := u.tagRepository.GetListByTagIdAndUserIdList(tagList)
	if err != nil {
		return nil, err
	}

	return tags, nil
}
