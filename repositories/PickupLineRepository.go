package repositories

import (
	"playbook/database"
	"playbook/entities"

	"github.com/google/uuid"
)

type PickupLineRepositoryInterface interface {
	Create(pickupLine *entities.PickupLine) error
	Update(pickupLine *entities.PickupLine) error
	Delete(pickupLine *entities.PickupLine) error
	GetByIdAndUserId(pickupLineUUID uuid.UUID, userUUID uuid.UUID) (*entities.PickupLine, error)
	GetById(pickupLineUUID uuid.UUID) (*entities.PickupLine, error)
	GetReactionByPickupLineIdAndUserId(pickupLineUUID uuid.UUID, userUUID uuid.UUID) (*entities.Reaction, error)
	CreateReaction(reaction *entities.Reaction) error
	UpdateReaction(reaction *entities.Reaction) error
	DeleteReaction(reaction *entities.Reaction) error
	GetStatisticsById(pickupLineUUID uuid.UUID) (*entities.Statistic, error)
	GetPickupLineForVisibilityCheck(pickupLineUUID uuid.UUID) (*entities.PickupLine, error)
}

type PickupLineRepository struct {
	database *database.Database
}

func (u *PickupLineRepository) Create(pickupLine *entities.PickupLine) error {
	return u.database.CreatePickupLine(pickupLine)
}

func (u *PickupLineRepository) Update(pickupLine *entities.PickupLine) error {
	return u.database.UpdatePickupLine(pickupLine)
}

func (u *PickupLineRepository) Delete(pickupLine *entities.PickupLine) error {
	return u.database.DeletePickupLine(pickupLine)
}

func CreatePickupLineRepository(database *database.Database) *PickupLineRepository {
	return &PickupLineRepository{database: database}
}

func (u *PickupLineRepository) GetByIdAndUserId(pickupLineId uuid.UUID, userId uuid.UUID) (*entities.PickupLine, error) {
	return u.database.GetPickupLineByIdAndUserId(pickupLineId, userId)
}

func (u *PickupLineRepository) GetById(pickupLineId uuid.UUID) (*entities.PickupLine, error) {
	return u.database.GetPickupLineById(pickupLineId)
}

func (u *PickupLineRepository) GetReactionByPickupLineIdAndUserId(pickupLineId uuid.UUID, userId uuid.UUID) (*entities.Reaction, error) {
	return u.database.GetReactionByPickupLineIdAndUserId(pickupLineId, userId)
}

func (u *PickupLineRepository) CreateReaction(reaction *entities.Reaction) error {
	return u.database.CreateReaction(reaction)
}

func (u *PickupLineRepository) UpdateReaction(reaction *entities.Reaction) error {
	if reaction.Vote == "" {
		reaction.Vote = entities.None
	}
	return u.database.UpdateReaction(reaction)
}

func (u *PickupLineRepository) DeleteReaction(reaction *entities.Reaction) error {
	return u.database.DeleteReaction(reaction)
}

func (u *PickupLineRepository) GetStatisticsById(pickupLineId uuid.UUID) (*entities.Statistic, error) {
	return u.database.GetStatisticsById(pickupLineId)
}

func (u *PickupLineRepository) GetPickupLineForVisibilityCheck(pickupLineUUID uuid.UUID) (*entities.PickupLine, error) {
	return u.database.GetPickupLineForVisibilityCheck(pickupLineUUID)
}
