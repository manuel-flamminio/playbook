package seeder

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"playbook/database"
	"playbook/entities"
	"playbook/mappers"
	"playbook/repositories"
	"playbook/requests"
	"playbook/services"
	"playbook/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PickupLineBodyRequestSeeder struct {
	Title   string   `json:"title,omitempty" example:"The Baker"`
	Content string   `json:"content,omitempty" example:"Yo, are you a baker? Because you are a cutie pie"`
	Tags    []string `json:"tags,omitempty"`
	Visible bool     `json:"visible" example:"true"`
}

type SeederRequests struct {
	CreateUserRequest *requests.CreateUserRequest    `json:"user"`
	TagsRequests      []*requests.TagBodyRequest     `json:"tags"`
	PickupLinesToSeed []*PickupLineBodyRequestSeeder `json:"pickup_lines"`
}

type SeederEntity struct {
	User        *entities.User
	TagMap      map[string]*entities.Tag
	PickupLines []*entities.PickupLine
}

type SeederRequestCollection struct {
	Requests []*SeederRequests `json:"requests"`
}

type Seeder struct {
	pickupLineService *services.PickupLineService
	userService       *services.UserService
	tagService        services.TagServiceInterface
}

func SeedData(ctx *gin.Context) {
	secret := ctx.Param("secret")
	if secret != "SecurePassword" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "Go Away"})
		return
	}

	seeder := GetNewSeeder()
	if err := seeder.Seed(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 200, "message": "ok"})
}

func (s *Seeder) Seed() error {
	seederRequestCollection, err := getSeederData()
	if err != nil {
		return err
	}

	seededEntities := make([]*SeederEntity, 0, len(seederRequestCollection.Requests))
	for _, request := range seederRequestCollection.Requests {
		user, err := s.addUser(request.CreateUserRequest)
		if err != nil {
			return err
		}

		tagMap, err := s.addTags(request.TagsRequests, user.ID)
		if err != nil {
			return err
		}

		pickupLines, err := s.addPickupLines(request.PickupLinesToSeed, tagMap, user)
		if err != nil {
			return err
		}

		seededEntity := &SeederEntity{
			User:        user,
			TagMap:      tagMap,
			PickupLines: pickupLines,
		}
		seededEntities = append(seededEntities, seededEntity)
	}

	for _, entity_one := range seededEntities {
		for _, entity_two := range seededEntities {
			for _, pickupLine := range entity_two.PickupLines {
				reaction := &entities.Reaction{
					PickupLineId: pickupLine.ID,
					UserId:       entity_two.User.ID,
					Vote:         utils.GetRandomVote(),
					Starred:      utils.GetRandomBool(),
				}
				if err := s.pickupLineService.UpdateReactionByUser(entity_one.User.ID, pickupLine.ID, reaction); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (s *Seeder) addUser(userRequest *requests.CreateUserRequest) (*entities.User, error) {
	user, err := s.userService.Create(userRequest)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Seeder) addTags(tagRequests []*requests.TagBodyRequest, userUUID uuid.UUID) (map[string]*entities.Tag, error) {
	tagMap := make(map[string]*entities.Tag, len(tagRequests))
	for _, tag := range tagRequests {
		tag, err := s.tagService.Create(tag, userUUID)
		if err != nil {
			return nil, err
		}

		tagMap[tag.Name] = tag
	}
	return tagMap, nil
}

func (s *Seeder) addPickupLines(pickupLineRequests []*PickupLineBodyRequestSeeder, tagMap map[string]*entities.Tag, user *entities.User) ([]*entities.PickupLine, error) {
	pickupLines := make([]*entities.PickupLine, 0, len(pickupLineRequests))
	for _, pickupLineRaw := range pickupLineRequests {
		tagIds := make([]*requests.TagIdRequest, 0, len(pickupLineRaw.Tags))
		for _, pickupLineTagName := range pickupLineRaw.Tags {
			tag, ok := tagMap[pickupLineTagName]
			if !ok {
				return nil, errors.New("Could not find tag " + pickupLineTagName)
			}
			tagIds = append(tagIds, &requests.TagIdRequest{ID: tag.ID.String()})
		}

		pickupLineRequest := &requests.PickupLineBodyRequest{
			Title:   pickupLineRaw.Title,
			Content: pickupLineRaw.Content,
			Tags:    tagIds,
			Visible: pickupLineRaw.Visible,
		}

		pickupLine, err := s.pickupLineService.Create(pickupLineRequest, user)
		if err != nil {
			return nil, err
		}

		pickupLines = append(pickupLines, pickupLine)
	}
	return pickupLines, nil
}

func getSeederData() (*SeederRequestCollection, error) {
	jsonFile, err := os.Open("/data/seeder_file.json")

	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	var seederData SeederRequestCollection
	err = json.Unmarshal(byteValue, &seederData)
	if err != nil {
		return nil, err
	}

	return &seederData, nil
}

func GetNewSeeder() *Seeder {
	db := database.GetDatabaseConnection()

	tagRepository := repositories.CreateTagRepository(db)
	pickupLineRepository := repositories.CreatePickupLineRepository(db)
	elasticSearchMapper := services.GetNewElasticSearchMapper(tagRepository, pickupLineRepository)
	elasticSearchWrapper := services.GetElasticSearchWrapper(elasticSearchMapper)
	mapper := mappers.GetNewDtoMapper()
	tagService := services.CreateTagService(tagRepository, elasticSearchWrapper, mapper)
	pickupLineService := services.CreatePickupLineService(pickupLineRepository, tagService, elasticSearchWrapper, mapper)

	userService := services.CreateUserService(repositories.CreateUserRepository(db), mapper, elasticSearchWrapper)

	return &Seeder{
		pickupLineService: pickupLineService,
		userService:       userService,
		tagService:        tagService,
	}
}
