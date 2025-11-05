package services

import (
	"playbook/constants"
	"playbook/entities"
	"playbook/filters"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
	"github.com/google/uuid"
)

const SORT_BY_NEW_SCALE_IN_DAYS = 1
const SORT_BY_TRENDING_SCALE_IN_DAYS = 1
const SORT_BY_TRENDING_UPVOTE_WEIGHT = 1.5
const SORT_BY_BEST_OF_ALL_TIME_UPVOTE_WEIGHT = 1.5

type QueryWrapper struct {
	query *types.Query
}

type QueryWrapperInterface interface {
	WithUserFilter(userUUID uuid.UUID)
	WithVisibleFilter(visible bool)
	WithStarredByUserFilter(userUUID uuid.UUID)
	WithUpvotedByUserFilter(userUUID uuid.UUID)
	WithSuccessPercentageFilter(successPercentage float64)
	WithTitleSuggestionShould(title string)
	WithContentSuggestion(content string)
	WithTagsShould(tagIds []string)
	WithUsernameSuggestionShould(username string)
	WithDisplayNameSuggestionShould(displayName string)
	GetQuery() *types.Query
	SetMinimumShouldMatch(minimumShouldMatch int)
	WithTimeScoring(scaleInDays int)
	WithUpvotesScoring(scoringWeight float64)
	WithRandomScoring()
}

type SearchRequestWrapper struct {
	request *search.Request
	query   QueryWrapperInterface
}

type SearchRequestWrapperInterface interface {
	withOffset(offset int)
	withSize(size int)
	applyFilters(filters *filters.PickupLineQueryFilters, requestingUser uuid.UUID)
	applySorting(pickupLineFilters *filters.PickupLineQueryFilters)
	applyUserFilters(filters *filters.UserQueryFilters)
	getRequest() *search.Request
}

func GetNewQueryWrapper() QueryWrapperInterface {
	query := &types.Query{
		Bool: &types.BoolQuery{},
		FunctionScore: &types.FunctionScoreQuery{
			Functions: make([]types.FunctionScore, 0, 0),
			Query:     &types.Query{},
		},
	}
	return &QueryWrapper{query: query}
}

func GetNewElasticSearchRequestWrapper(queryWrapper QueryWrapperInterface) SearchRequestWrapperInterface {
	return &SearchRequestWrapper{
		request: &search.Request{},
		query:   queryWrapper,
	}
}

func (e *SearchRequestWrapper) withOffset(offset int) {
	e.request.From = &offset
}

func (e *SearchRequestWrapper) withSize(size int) {
	e.request.Size = &size
}

func (q *QueryWrapper) WithUserFilter(userUUID uuid.UUID) {
	userQuery := &types.Query{
		Term: map[string]types.TermQuery{
			constants.ELASTIC_USER_ID_FIELD: {Value: userUUID},
		},
	}

	q.addToFilterQuery(userQuery)
}

func (q *QueryWrapper) WithTimeScoring(scaleInDays int) {
	timeGaussFunction := &types.FunctionScore{
		Gauss: &types.DecayFunctionBaseDateMathDuration{
			DecayFunctionBaseDateMathDuration: map[string]types.DecayPlacementDateMathDuration{
				constants.ELASTIC_UPDATED_AT_FIELD: types.DecayPlacementDateMathDuration{
					Scale: strconv.Itoa(scaleInDays) + "d",
				},
			},
		},
	}

	q.query.FunctionScore.Functions = append(q.query.FunctionScore.Functions, *timeGaussFunction)
}

func (q *QueryWrapper) WithUpvotesScoring(scoringWeight float64) {
	upvoteFactorFunction := &types.FunctionScore{
		FieldValueFactor: &types.FieldValueFactorScoreFunction{
			Factor: (*types.Float64)(&scoringWeight),
			Field:  constants.ELASTIC_NUMBER_OF_SUCCESSES,
		},
	}

	q.query.FunctionScore.Functions = append(q.query.FunctionScore.Functions, *upvoteFactorFunction)
}

func (q *QueryWrapper) WithRandomScoring() {
	randomScoringField := "_seq_no"
	randomScoreFunction := &types.FunctionScore{
		RandomScore: &types.RandomScoreFunction{
			Field: &randomScoringField,
		},
	}

	q.query.FunctionScore.Functions = append(q.query.FunctionScore.Functions, *randomScoreFunction)
}

func (q *QueryWrapper) SetMinimumShouldMatch(minimumShouldMatch int) {
	q.query.Bool.MinimumShouldMatch = minimumShouldMatch
}

func (q *QueryWrapper) WithVisibleFilter(visible bool) {
	visibleQuery := &types.Query{
		Term: map[string]types.TermQuery{
			constants.ELASTIC_VISIBLE_FIELD: {Value: visible},
		},
	}
	q.addToFilterQuery(visibleQuery)
}

func (q *QueryWrapper) WithStarredByUserFilter(userUUID uuid.UUID) {
	starredQuery := &types.Query{
		Term: map[string]types.TermQuery{
			constants.ELASTIC_STARRED_BY_USER_FIELD: {Value: userUUID},
		},
	}
	q.addToFilterQuery(starredQuery)
}

func (q *QueryWrapper) WithUpvotedByUserFilter(userUUID uuid.UUID) {
	upvotedQuery := &types.Query{
		Term: map[string]types.TermQuery{
			constants.ELASTIC_UPVOTED_BY_USER_FIELD: {Value: userUUID},
		},
	}
	q.addToFilterQuery(upvotedQuery)
}

func (q *QueryWrapper) WithSuccessPercentageFilter(successPercentage float64) {
	successPercentageQuery := &types.Query{
		Range: map[string]types.RangeQuery{
			constants.ELASTIC_SUCCESS_PERCENTAGE_FIELD: &types.NumberRangeQuery{Gte: (*types.Float64)(&successPercentage)},
		},
	}
	q.addToFilterQuery(successPercentageQuery)
}

func (q *QueryWrapper) WithAutocompleteSuggestionShould(fieldName string, value string) {
	autocompleteQuery := &types.Query{
		MultiMatch: &types.MultiMatchQuery{
			Query: value,
			Type:  &textquerytype.Boolprefix,
			Fields: []string{
				fieldName,
				fieldName + ".2gram",
				fieldName + ".3gram",
			},
		},
	}
	q.addToShouldQuery(autocompleteQuery)
}

func (q *QueryWrapper) WithDisplayNameSuggestionShould(displayName string) {
	q.WithAutocompleteSuggestionShould(constants.ELASTIC_DISPLAY_NAME_FIELD, displayName)
}

func (q *QueryWrapper) WithUsernameSuggestionShould(username string) {
	q.WithAutocompleteSuggestionShould(constants.ELASTIC_USERNAME_FIELD, username)
}

func (q *QueryWrapper) WithTitleSuggestionShould(title string) {
	q.WithAutocompleteSuggestionShould(constants.ELASTIC_TITLE_FIELD, title)
}

func (q *QueryWrapper) WithContentSuggestion(content string) {
	contentQuery := &types.Query{
		Match: map[string]types.MatchQuery{
			constants.ELASTIC_CONTENT_FIELD: {Query: content},
		},
	}
	q.addToShouldQuery(contentQuery)
}

func (q *QueryWrapper) WithTagsShould(tagIds []string) {
	for _, tagId := range tagIds {
		tagQuery := &types.Query{
			Term: map[string]types.TermQuery{
				constants.ELASTIC_TAG_FIELD: {Value: tagId},
			},
		}
		q.addToShouldQuery(tagQuery)
	}
}

func (q *QueryWrapper) GetQuery() *types.Query {
	if len(q.query.FunctionScore.Functions) > 0 {
		q.query.FunctionScore.Query.Bool = q.query.Bool
		q.query.FunctionScore.Query.Bool.Must = q.query.Bool.Filter
		q.query.FunctionScore.Query.Bool.Filter = nil
		q.query.Bool = nil
		return q.query
	}

	q.query.FunctionScore = nil
	return q.query
}

func (e *SearchRequestWrapper) applyUserFilters(filters *filters.UserQueryFilters) {
	e.withOffset(filters.Page * constants.ITEMS_FOR_PAGE)

	if filters.Username != "" {
		e.query.WithUsernameSuggestionShould(filters.Username)
	}

	if filters.DisplayName != "" {
		e.query.WithDisplayNameSuggestionShould(filters.DisplayName)
	}
}

func (e *SearchRequestWrapper) applySorting(pickupLineFilters *filters.PickupLineQueryFilters) {
	switch pickupLineFilters.SortingType {
	case filters.New:
		e.query.WithTimeScoring(SORT_BY_NEW_SCALE_IN_DAYS)
	case filters.BestOfAllTime:
		e.query.WithUpvotesScoring(SORT_BY_BEST_OF_ALL_TIME_UPVOTE_WEIGHT)
	case filters.Trending:
		e.query.WithTimeScoring(SORT_BY_TRENDING_SCALE_IN_DAYS)
		e.query.WithUpvotesScoring(SORT_BY_TRENDING_UPVOTE_WEIGHT)
	case filters.Random:
		e.query.WithRandomScoring()
	}
}

func (e *SearchRequestWrapper) applyFilters(filters *filters.PickupLineQueryFilters, requestingUser uuid.UUID) {
	e.withOffset(filters.Page * constants.ITEMS_FOR_PAGE)

	if len(filters.Tags) > 0 {
		e.query.WithTagsShould(filters.Tags)
	}

	if filters.Title != "" {
		e.query.WithTitleSuggestionShould(filters.Title)
	}

	if filters.Content != "" {
		e.query.WithContentSuggestion(filters.Content)
	}

	if filters.Starred {
		e.query.WithStarredByUserFilter(requestingUser)
	}

	if filters.OnlyUpvoted && (filters.UserId == uuid.Nil || filters.UserId == requestingUser) {
		e.query.WithUpvotedByUserFilter(requestingUser)
	}

	if filters.UserId != uuid.Nil {
		e.query.WithUserFilter(filters.UserId)
	}

	if filters.SuccessPercentage > 0.0 {
		e.query.WithSuccessPercentageFilter(filters.SuccessPercentage)
	}

	isRequestingSelfPickupLines := filters.UserId == requestingUser
	if filters.Visible != "" {
		switch filters.Visible {
		case entities.NotVisible:
			if isRequestingSelfPickupLines {
				e.query.WithVisibleFilter(false)
			} else {
				e.query.WithVisibleFilter(true)
			}
		case entities.Visible:
			e.query.WithVisibleFilter(true)
		default:
		}
	} else if !isRequestingSelfPickupLines {
		e.query.WithVisibleFilter(true)
	}
}

func (e *SearchRequestWrapper) getRequest() *search.Request {
	e.request.Query = e.query.GetQuery()
	return e.request
}

func (q *QueryWrapper) addToFilterQuery(query *types.Query) {
	q.query.Bool.Filter = append(q.query.Bool.Filter, *query)
}

func (q *QueryWrapper) addToShouldQuery(query *types.Query) {
	q.query.Bool.Should = append(q.query.Bool.Should, *query)
}
