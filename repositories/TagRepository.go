package repositories

import (
	"playbook/database"
	"playbook/entities"

	"github.com/google/uuid"
)

type TagRepositoryInterface interface {
	Create(tag *entities.Tag) error
	Update(tag *entities.Tag) error
	Delete(tag *entities.Tag) error
	GetById(tagId string) (*entities.Tag, error)
	GetByUserId(userId uuid.UUID) ([]*entities.Tag, error)
	GetListByTagIdAndUserIdList(tagList [][]uuid.UUID) ([]*entities.Tag, error)
}

type TagRepository struct {
	database *database.Database
}

func (u *TagRepository) Create(tag *entities.Tag) error {
	return u.database.CreateTag(tag)
}

func (u *TagRepository) Update(tag *entities.Tag) error {
	return u.database.UpdateTag(tag)
}

func (u *TagRepository) Delete(tag *entities.Tag) error {
	return u.database.DeleteTag(tag)
}

func CreateTagRepository(database *database.Database) *TagRepository {
	return &TagRepository{database: database}
}

func (u *TagRepository) GetById(tagId string) (*entities.Tag, error) {
	return u.database.GetTagById(tagId)
}

func (u *TagRepository) GetByUserId(userId uuid.UUID) ([]*entities.Tag, error) {
	return u.database.GetTagsByUserId(userId)
}

func (u *TagRepository) GetListByTagIdAndUserIdList(tagList [][]uuid.UUID) ([]*entities.Tag, error) {
	return u.database.GetListByTagIdAndUserIdList(tagList)
}
