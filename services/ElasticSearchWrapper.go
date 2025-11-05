package services

import (
	"context"
	"encoding/json"
	"os"
	"playbook/constants"
	"playbook/entities"
	"playbook/filters"
	"playbook/responses"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/deletebyquery"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/indices/create"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/google/uuid"
)

type ElasticSearchWrapperInterface interface {
	ExistsIndex(indexName string) (bool, error)
	CreatePickupLineIndex() error
	CreateTagIndex() error
	CreateUserIndex() error
	IndexPickupLine(pickupLine *entities.PickupLine) error
	DeletePickupLine(pickupLineUUID uuid.UUID) error
	DeleteUserPickupLines(userUUID uuid.UUID) error
	DeleteTag(tag *entities.Tag) error
	IndexTag(tag *entities.Tag) (string, error)
	SearchPickupLines(userId uuid.UUID, filters *filters.PickupLineQueryFilters) (*responses.ElasticSearchPickupLineResponse, error)
	GetPickupLineFeed(userId uuid.UUID, filters *filters.PickupLineQueryFilters) (*responses.ElasticSearchPickupLineResponse, error)
	UpdatePickupLine(pickupLine *entities.PickupLine) error
	UpdateUserReaction(userId uuid.UUID, pickupLineUUID uuid.UUID, newReaction *entities.Reaction, oldReaction *entities.Reaction) error
	UpdateTag(tag *entities.Tag) error
	IndexUser(user *entities.User) error
	DeleteUser(userUUID uuid.UUID) error
	SearchUsers(filters *filters.UserQueryFilters) (*responses.ElasticSearchUserListResponse, error)
}

type ElasticSearchWrapper struct {
	client              *elasticsearch.TypedClient
	mapper              ElasticSearchMapperInterface
	PickupLineIndexName string
	TagsIndexName       string
	UserIndexName       string
}

type TagName struct {
	Name []string `json:"name"`
}

func getClient() *elasticsearch.TypedClient {
	cert, err := os.ReadFile(os.Getenv(constants.ELASTIC_SEARCH_CERT_PATH))
	if err != nil {
		panic(err.Error())
	}

	cfg := elasticsearch.Config{
		Addresses: []string{
			os.Getenv(constants.ELASTIC_SEARCH_HOST),
		},
		Username: os.Getenv(constants.ELASTIC_SEARCH_USERNAME),
		Password: os.Getenv(constants.ELASTIC_SEARCH_PASSWORD),
		CACert:   cert,
	}

	es, err := elasticsearch.NewTypedClient(cfg)
	if err != nil {
		panic("Could not connect to elasticSearch")
	}
	return es
}

func GetElasticSearchWrapper(mapper ElasticSearchMapperInterface) *ElasticSearchWrapper {
	client := getClient()
	return &ElasticSearchWrapper{
		client:              client,
		mapper:              mapper,
		PickupLineIndexName: os.Getenv(constants.PICKUP_LINE_INDEX_NAME),
		TagsIndexName:       os.Getenv(constants.TAGS_INDEX_NAME),
		UserIndexName:       os.Getenv(constants.USER_INDEX_NAME),
	}
}

func (e *ElasticSearchWrapper) ExistsIndex(indexName string) (bool, error) {
	return e.client.Indices.Exists(indexName).Do(context.Background())
}

func (e *ElasticSearchWrapper) CreateUserIndex() error {
	searchAsYouTypeProperty := types.NewSearchAsYouTypeProperty()
	MaxShingleSize := 3
	searchAsYouTypeProperty.MaxShingleSize = &MaxShingleSize

	_, err := e.client.Indices.Create(e.UserIndexName).Request(
		&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					constants.ELASTIC_ID_FIELD:           types.NewKeywordProperty(),
					constants.ELASTIC_DISPLAY_NAME_FIELD: searchAsYouTypeProperty,
					constants.ELASTIC_USERNAME_FIELD:     searchAsYouTypeProperty,
				},
			},
		},
	).Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (e *ElasticSearchWrapper) CreatePickupLineIndex() error {
	searchAsYouTypeProperty := types.NewSearchAsYouTypeProperty()
	MaxShingleSize := 3
	searchAsYouTypeProperty.MaxShingleSize = &MaxShingleSize

	_, err := e.client.Indices.Create(e.PickupLineIndexName).Request(
		&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					constants.ELASTIC_ID_FIELD:                 types.NewKeywordProperty(),
					constants.ELASTIC_TITLE_FIELD:              searchAsYouTypeProperty,
					constants.ELASTIC_CONTENT_FIELD:            types.NewTextProperty(),
					constants.ELASTIC_NUMBER_OF_TRIES:          types.NewUnsignedLongNumberProperty(),
					constants.ELASTIC_NUMBER_OF_FAILURES:       types.NewUnsignedLongNumberProperty(),
					constants.ELASTIC_NUMBER_OF_SUCCESSES:      types.NewUnsignedLongNumberProperty(),
					constants.ELASTIC_SUCCESS_PERCENTAGE_FIELD: types.NewFloatNumberProperty(),
					constants.ELASTIC_STARRED_FIELD:            types.NewBooleanProperty(),
					constants.ELASTIC_STARRED_BY_USER_FIELD:    types.NewKeywordProperty(),
					constants.ELASTIC_UPVOTED_BY_USER_FIELD:    types.NewKeywordProperty(),
					constants.ELASTIC_DOWNVOTED_BY_USER_FIELD:  types.NewKeywordProperty(),
					constants.ELASTIC_TAG_FIELD:                types.NewKeywordProperty(),
					constants.ELASTIC_VISIBLE_FIELD:            types.NewBooleanProperty(),
					constants.ELASTIC_DISPLAY_NAME_FIELD:       types.NewKeywordProperty(),
					constants.ELASTIC_USER_ID_FIELD:            types.NewKeywordProperty(),
					constants.ELASTIC_USERNAME_FIELD:           types.NewKeywordProperty(),
					constants.ELASTIC_UPDATED_AT_FIELD:         types.NewDateProperty(),
				},
			},
		},
	).Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (e *ElasticSearchWrapper) CreateTagIndex() error {
	_, err := e.client.Indices.Create(e.TagsIndexName).Request(
		&create.Request{
			Mappings: &types.TypeMapping{
				Properties: map[string]types.Property{
					"id":     types.NewKeywordProperty(),
					"name":   types.NewKeywordProperty(),
					"userId": types.NewKeywordProperty(),
				},
			},
		},
	).Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (e *ElasticSearchWrapper) IndexPickupLine(pickupLine *entities.PickupLine) error {
	_, err := e.client.Index(e.PickupLineIndexName).Id(pickupLine.ID.String()).Request(e.mapper.pickupLineToElasticDTO(pickupLine)).Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (e *ElasticSearchWrapper) DeletePickupLine(pickupLineUUID uuid.UUID) error {
	_, err := e.client.Delete(e.PickupLineIndexName, pickupLineUUID.String()).Do(context.Background())
	return err
}

func (e *ElasticSearchWrapper) DeleteUserPickupLines(userUUID uuid.UUID) error {
	query := GetNewQueryWrapper()
	query.WithUserFilter(userUUID)
	_, err := e.client.DeleteByQuery(e.PickupLineIndexName).Request(&deletebyquery.Request{Query: query.GetQuery()}).Do(context.Background())
	return err
}

func (e *ElasticSearchWrapper) DeleteTag(tag *entities.Tag) error {
	_, err := e.client.Delete(e.TagsIndexName, tag.ElasticSearchId).Do(context.Background())
	return err
}

func (e *ElasticSearchWrapper) DeleteUser(userUUID uuid.UUID) error {
	_, err := e.client.Delete(e.UserIndexName, userUUID.String()).Do(context.Background())
	return err
}

func (e *ElasticSearchWrapper) IndexTag(tag *entities.Tag) (string, error) {
	res, err := e.client.Index(e.TagsIndexName).Request(e.mapper.tagToElasticDTO(tag)).Do(context.Background())
	if err != nil {
		return "", err
	}

	return res.Id_, nil
}

func (e *ElasticSearchWrapper) IndexUser(user *entities.User) error {
	_, err := e.client.Index(e.UserIndexName).Id(user.ID.String()).Request(e.mapper.userToElasticDTO(user)).Do(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (e *ElasticSearchWrapper) SearchUsers(filters *filters.UserQueryFilters) (*responses.ElasticSearchUserListResponse, error) {
	elasticSearchRequest := GetNewElasticSearchRequestWrapper(GetNewQueryWrapper())
	elasticSearchRequest.withSize(constants.ITEMS_FOR_PAGE)
	elasticSearchRequest.applyUserFilters(filters)

	res, err := e.client.Search().Index(e.UserIndexName).Request(elasticSearchRequest.getRequest()).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return e.mapper.hydrateUserListFromElasticSearchResponse(res)
}

func (e *ElasticSearchWrapper) SearchPickupLines(userId uuid.UUID, filters *filters.PickupLineQueryFilters) (*responses.ElasticSearchPickupLineResponse, error) {
	elasticSearchRequest := GetNewElasticSearchRequestWrapper(GetNewQueryWrapper())
	elasticSearchRequest.withSize(constants.ITEMS_FOR_PAGE)
	elasticSearchRequest.applyFilters(filters, userId)
	elasticSearchRequest.applySorting(filters)

	res, err := e.client.Search().Index(e.PickupLineIndexName).Request(elasticSearchRequest.getRequest()).Do(context.Background())
	if err != nil {
		return nil, err
	}

	return e.mapper.hydratePickupLineListFromElasticSearchResponse(res, userId)
}

func (e *ElasticSearchWrapper) GetPickupLineFeed(userId uuid.UUID, filters *filters.PickupLineQueryFilters) (*responses.ElasticSearchPickupLineResponse, error) {
	elasticSearchRequest := GetNewElasticSearchRequestWrapper(GetNewQueryWrapper())
	filters.Visible = entities.Visible
	elasticSearchRequest.withSize(constants.ITEMS_FOR_PAGE)
	elasticSearchRequest.applyFilters(filters, userId)
	elasticSearchRequest.applySorting(filters)

	res, err := e.client.Search().Index(e.PickupLineIndexName).Request(elasticSearchRequest.getRequest()).Do(context.Background())

	if err != nil {
		return nil, err
	}

	return e.mapper.hydratePickupLineListFromElasticSearchResponse(res, userId)
}

func (e *ElasticSearchWrapper) UpdatePickupLine(pickupLine *entities.PickupLine) error {
	pickupLineElasticDTO := e.mapper.pickupLineToUpdateElasticDTO(pickupLine)
	requestBody, err := json.Marshal(pickupLineElasticDTO)
	if err != nil {
		return err
	}

	request := &update.Request{
		Doc: requestBody,
	}
	_, err = e.client.Update(e.PickupLineIndexName, pickupLine.ID.String()).Request(request).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (e *ElasticSearchWrapper) UpdateUserReaction(userId uuid.UUID, pickupLineUUID uuid.UUID, newReaction *entities.Reaction, oldReaction *entities.Reaction) error {
	scriptWrapper := NewScriptWrapper()
	if oldReaction.Starred != newReaction.Starred {
		if newReaction.Starred {
			scriptWrapper.AddStarredByUser(userId)
		} else {
			scriptWrapper.RemoveStarredByUser(userId)
		}
	}

	if oldReaction.Vote != newReaction.Vote {
		switch newReaction.Vote {
		case entities.Upvote:
			scriptWrapper.AddUpvotedByUser(userId)
		case entities.Downvote:
			scriptWrapper.AddDownvotedByUser(userId)
		}

		switch oldReaction.Vote {
		case entities.Upvote:
			scriptWrapper.RemoveUpvotedByUser(userId)
		case entities.Downvote:
			scriptWrapper.RemoveDownvotedByUser(userId)
		}
	}

	script, err := scriptWrapper.GetScript()
	if err != nil {
		return err
	}

	request := &update.Request{
		Script: script,
	}
	_, err = e.client.Update(e.PickupLineIndexName, pickupLineUUID.String()).Request(request).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func (e *ElasticSearchWrapper) UpdateTag(tag *entities.Tag) error {
	tagElasticDTO := e.mapper.tagToElasticDTO(tag)
	requestBody, err := json.Marshal(tagElasticDTO)
	if err != nil {
		return err
	}

	request := &update.Request{
		Doc: requestBody,
	}
	_, err = e.client.Update(e.TagsIndexName, tag.ElasticSearchId).Request(request).Do(context.Background())
	if err != nil {
		return err
	}
	return nil
}
