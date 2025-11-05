package services

import (
	"math/rand"
	"playbook/constants"
	"playbook/entities"
	"playbook/filters"
	"playbook/mocks"
	"playbook/utils"
	"reflect"
	"strconv"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/textquerytype"
	"go.uber.org/mock/gomock"
)

func TestRandomScoring(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randomScoringField := "_seq_no"

	query := GetNewQueryWrapper()
	query.WithRandomScoring()
	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	randomScoreFunction := &types.FunctionScore{
		RandomScore: &types.RandomScoreFunction{
			Field: &randomScoringField,
		},
	}

	expectedQuery.query.FunctionScore.Functions = append(expectedQuery.query.FunctionScore.Functions, *randomScoreFunction)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestUpvotesScoring(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scoringWeight := 1.22

	query := GetNewQueryWrapper()
	query.WithUpvotesScoring(scoringWeight)
	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	upvoteFactorFunction := &types.FunctionScore{
		FieldValueFactor: &types.FieldValueFactorScoreFunction{
			Factor: (*types.Float64)(&scoringWeight),
			Field:  constants.ELASTIC_NUMBER_OF_SUCCESSES,
		},
	}

	expectedQuery.query.FunctionScore.Functions = append(expectedQuery.query.FunctionScore.Functions, *upvoteFactorFunction)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestTimeScoring(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	scaleInDays := 1

	query := GetNewQueryWrapper()
	query.WithTimeScoring(scaleInDays)
	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	timeGaussFunction := &types.FunctionScore{
		Gauss: &types.DecayFunctionBaseDateMathDuration{
			DecayFunctionBaseDateMathDuration: map[string]types.DecayPlacementDateMathDuration{
				constants.ELASTIC_UPDATED_AT_FIELD: types.DecayPlacementDateMathDuration{
					Scale: strconv.Itoa(scaleInDays) + "d",
				},
			},
		},
	}

	expectedQuery.query.FunctionScore.Functions = append(expectedQuery.query.FunctionScore.Functions, *timeGaussFunction)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestWithUserFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()

	query := GetNewQueryWrapper()
	query.WithUserFilter(userUUID)
	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	userQuery := &types.Query{
		Term: map[string]types.TermQuery{
			constants.ELASTIC_USER_ID_FIELD: {Value: userUUID},
		},
	}
	expectedQuery.query.Bool.Filter = append(expectedQuery.query.Bool.Filter, *userQuery)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestSetMinimumShouldMatch(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randInt := rand.Int()
	query := GetNewQueryWrapper()
	query.SetMinimumShouldMatch(randInt)
	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	expectedQuery.query.Bool.MinimumShouldMatch = randInt

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestWithVisibleFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	randBool := utils.GetRandomBool()
	query := GetNewQueryWrapper()
	query.WithVisibleFilter(randBool)
	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	visibleQuery := &types.Query{
		Term: map[string]types.TermQuery{
			constants.ELASTIC_VISIBLE_FIELD: {Value: randBool},
		},
	}
	expectedQuery.query.Bool.Filter = append(expectedQuery.query.Bool.Filter, *visibleQuery)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestWithOnlyUpvotedFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	query := GetNewQueryWrapper()
	query.WithUpvotedByUserFilter(userUUID)
	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	upvotedQuery := &types.Query{
		Term: map[string]types.TermQuery{
			constants.ELASTIC_UPVOTED_BY_USER_FIELD: {Value: userUUID},
		},
	}
	expectedQuery.query.Bool.Filter = append(expectedQuery.query.Bool.Filter, *upvotedQuery)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestWithStarredFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	query := GetNewQueryWrapper()
	query.WithStarredByUserFilter(userUUID)
	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	starredQuery := &types.Query{
		Term: map[string]types.TermQuery{
			constants.ELASTIC_STARRED_BY_USER_FIELD: {Value: userUUID},
		},
	}
	expectedQuery.query.Bool.Filter = append(expectedQuery.query.Bool.Filter, *starredQuery)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestWithSuccessPercentageFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	successPercentage := rand.Float64()
	query := GetNewQueryWrapper()
	query.WithSuccessPercentageFilter(successPercentage)
	successPercentageQuery := &types.Query{
		Range: map[string]types.RangeQuery{
			constants.ELASTIC_SUCCESS_PERCENTAGE_FIELD: &types.NumberRangeQuery{Gte: (*types.Float64)(&successPercentage)},
		},
	}
	expectedQuery.query.Bool.Filter = append(expectedQuery.query.Bool.Filter, *successPercentageQuery)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestWithTitleSuggestionShould(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	title := utils.GetRandomString(10)
	query := GetNewQueryWrapper()
	query.WithTitleSuggestionShould(title)
	titleQuery := &types.Query{
		MultiMatch: &types.MultiMatchQuery{
			Query: title,
			Type:  &textquerytype.Boolprefix,
			Fields: []string{
				constants.ELASTIC_TITLE_FIELD,
				constants.ELASTIC_TITLE_FIELD + ".2gram",
				constants.ELASTIC_TITLE_FIELD + ".3gram",
			},
		},
	}
	expectedQuery.query.Bool.Should = append(expectedQuery.query.Bool.Should, *titleQuery)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestWithDisplayNameSuggestionShould(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	displayName := utils.GetRandomString(10)
	query := GetNewQueryWrapper()
	query.WithDisplayNameSuggestionShould(displayName)
	displayNameQuery := &types.Query{
		MultiMatch: &types.MultiMatchQuery{
			Query: displayName,
			Type:  &textquerytype.Boolprefix,
			Fields: []string{
				constants.ELASTIC_DISPLAY_NAME_FIELD,
				constants.ELASTIC_DISPLAY_NAME_FIELD + ".2gram",
				constants.ELASTIC_DISPLAY_NAME_FIELD + ".3gram",
			},
		},
	}
	expectedQuery.query.Bool.Should = append(expectedQuery.query.Bool.Should, *displayNameQuery)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestWithUsernameSuggestionShould(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	username := utils.GetRandomUsername()
	query := GetNewQueryWrapper()
	query.WithUsernameSuggestionShould(username)
	usernameQuery := &types.Query{
		MultiMatch: &types.MultiMatchQuery{
			Query: username,
			Type:  &textquerytype.Boolprefix,
			Fields: []string{
				constants.ELASTIC_USERNAME_FIELD,
				constants.ELASTIC_USERNAME_FIELD + ".2gram",
				constants.ELASTIC_USERNAME_FIELD + ".3gram",
			},
		},
	}
	expectedQuery.query.Bool.Should = append(expectedQuery.query.Bool.Should, *usernameQuery)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestWithContentSuggestionShould(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	content := utils.GetRandomString(10)
	contentQuery := &types.Query{
		Match: map[string]types.MatchQuery{
			constants.ELASTIC_CONTENT_FIELD: {Query: content},
		},
	}
	query := GetNewQueryWrapper()
	query.WithContentSuggestion(content)
	expectedQuery.query.Bool.Should = append(expectedQuery.query.Bool.Should, *contentQuery)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestWithTagShould(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	numberOfTags := 10
	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	tagIds := make([]string, 0, numberOfTags)
	for i := 0; i < numberOfTags; i++ {
		randomUUID, _ := utils.GetRandomUUID()
		tagIds = append(tagIds, randomUUID.String())
		tagQuery := &types.Query{
			Term: map[string]types.TermQuery{
				constants.ELASTIC_TAG_FIELD: {Value: randomUUID.String()},
			},
		}
		expectedQuery.query.Bool.Should = append(expectedQuery.query.Bool.Should, *tagQuery)
	}
	query := GetNewQueryWrapper()
	query.WithTagsShould(tagIds)

	if query == nil || !reflect.DeepEqual(query, expectedQuery) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestGetQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	query := GetNewQueryWrapper()
	queryObj := query.GetQuery()
	expectedQuery.query.FunctionScore = nil

	if queryObj == nil || !reflect.DeepEqual(queryObj, expectedQuery.query) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestGetQueryWithScoringFunctions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	expectedQuery := GetNewQueryWrapper().(*QueryWrapper)
	query := GetNewQueryWrapper()
	query.WithRandomScoring()
	queryObj := query.GetQuery()

	randomScoringField := "_seq_no"
	randomScoreFunction := &types.FunctionScore{
		RandomScore: &types.RandomScoreFunction{
			Field: &randomScoringField,
		},
	}
	expectedQuery.query.FunctionScore.Functions = append(expectedQuery.query.FunctionScore.Functions, *randomScoreFunction)
	expectedQuery.query.FunctionScore.Query.Bool = expectedQuery.query.Bool
	expectedQuery.query.FunctionScore.Query.Bool.Must = expectedQuery.query.Bool.Filter
	expectedQuery.query.FunctionScore.Query.Bool.Filter = nil
	expectedQuery.query.Bool = nil

	if queryObj == nil || !reflect.DeepEqual(queryObj, expectedQuery.query) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestWithSize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().GetQuery().Return(nil)

	size := 10
	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	expectedsearchRequestWrapper.getRequest().Size = &size

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.withSize(size)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyUserFiltersPageFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filters := &filters.UserQueryFilters{
		Page: 2,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := filters.Page * constants.ITEMS_FOR_PAGE
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyUserFilters(filters)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyUserFiltersUsernameFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filters := &filters.UserQueryFilters{
		Page:     2,
		Username: utils.GetRandomUsername(),
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithUsernameSuggestionShould(filters.Username)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := filters.Page * constants.ITEMS_FOR_PAGE
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyUserFilters(filters)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyScoringBestOfAllTime(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filters := &filters.PickupLineQueryFilters{
		SortingType: filters.BestOfAllTime,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithUpvotesScoring(SORT_BY_BEST_OF_ALL_TIME_UPVOTE_WEIGHT)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applySorting(filters)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyScoringRandom(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filters := &filters.PickupLineQueryFilters{
		SortingType: filters.Random,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithRandomScoring()

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applySorting(filters)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyScoringNew(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filters := &filters.PickupLineQueryFilters{
		SortingType: filters.New,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithTimeScoring(SORT_BY_NEW_SCALE_IN_DAYS)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applySorting(filters)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyScoringTrending(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filters := &filters.PickupLineQueryFilters{
		SortingType: filters.Trending,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithTimeScoring(SORT_BY_TRENDING_SCALE_IN_DAYS)
	mockQuery.EXPECT().WithUpvotesScoring(SORT_BY_TRENDING_UPVOTE_WEIGHT)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applySorting(filters)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyFiltersOnlyUpvotedFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	filters := &filters.PickupLineQueryFilters{
		Page:        2,
		OnlyUpvoted: true,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().WithUpvotedByUserFilter(userUUID)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := filters.Page * constants.ITEMS_FOR_PAGE
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyFiltersOnlyUpvotedFilterWithOtherUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	otherUserUUID, _ := utils.GetRandomUUID()
	filters := &filters.PickupLineQueryFilters{
		Page:        2,
		OnlyUpvoted: true,
		UserId:      otherUserUUID,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().WithUserFilter(otherUserUUID)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := filters.Page * constants.ITEMS_FOR_PAGE
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyFiltersPageFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	filters := &filters.PickupLineQueryFilters{
		Page: 2,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := filters.Page * constants.ITEMS_FOR_PAGE
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyFiltersTagsFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	numberOfTags := 10
	userUUID, _ := utils.GetRandomUUID()
	tagIds := make([]string, 0, numberOfTags)
	for i := 0; i < numberOfTags; i++ {
		randomUUID, _ := utils.GetRandomUUID()
		tagIds = append(tagIds, randomUUID.String())
	}

	filters := &filters.PickupLineQueryFilters{
		Tags: tagIds,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithTagsShould(tagIds)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyFiltersTitleSuggestionShould(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	title := utils.GetRandomString(10)
	filters := &filters.PickupLineQueryFilters{
		Title: title,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithTitleSuggestionShould(title)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyFiltersContentSuggestionShould(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	content := utils.GetRandomString(10)
	filters := &filters.PickupLineQueryFilters{
		Content: content,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithContentSuggestion(content)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyStarredByUserFilters(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	filters := &filters.PickupLineQueryFilters{
		Starred: true,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithStarredByUserFilter(userUUID)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyUserIdSameAsRequestingFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	filters := &filters.PickupLineQueryFilters{
		UserId: userUUID,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithUserFilter(userUUID)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyUserIdDifferentFromRequestingFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	userToFilterUUID, _ := utils.GetRandomUUID()
	filters := &filters.PickupLineQueryFilters{
		UserId: userToFilterUUID,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithUserFilter(userToFilterUUID)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplySuccessPercentageFilter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	successPercentage := 0.33
	filters := &filters.PickupLineQueryFilters{
		SuccessPercentage: successPercentage,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithSuccessPercentageFilter(successPercentage)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyFiltersVisibilityAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	filters := &filters.PickupLineQueryFilters{
		Visible: entities.All,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyFiltersOnlyVisible(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	filters := &filters.PickupLineQueryFilters{
		Visible: entities.Visible,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyFiltersOnlyNotVisibleSameUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	filters := &filters.PickupLineQueryFilters{
		Visible: entities.NotVisible,
		UserId:  userUUID,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithUserFilter(userUUID)
	mockQuery.EXPECT().WithVisibleFilter(false)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyFiltersOnlyNotVisibleDifferentUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	userToFilterUUID, _ := utils.GetRandomUUID()
	filters := &filters.PickupLineQueryFilters{
		Visible: entities.NotVisible,
		UserId:  userToFilterUUID,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithUserFilter(userToFilterUUID)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}

func TestApplyFiltersOnlyNotVisibleNoUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userUUID, _ := utils.GetRandomUUID()
	filters := &filters.PickupLineQueryFilters{
		Visible: entities.NotVisible,
	}

	mockQuery := mocks.NewMockQueryWrapperInterface(ctrl)
	mockQuery.EXPECT().WithVisibleFilter(true)
	mockQuery.EXPECT().GetQuery().Return(nil)

	expectedsearchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	offset := 0
	expectedsearchRequestWrapper.getRequest().From = &offset

	searchRequestWrapper := GetNewElasticSearchRequestWrapper(mockQuery)
	searchRequestWrapper.applyFilters(filters, userUUID)
	if !reflect.DeepEqual(searchRequestWrapper, expectedsearchRequestWrapper) {
		t.Errorf("Error occurred in query wrapper")
		return
	}
}
