package controllers

import (
	"errors"
	"log"
	"net/http"
	"playbook/constants"
	"playbook/database"
	"playbook/entities"
	"playbook/filters"
	"playbook/mappers"
	"playbook/repositories"
	"playbook/requests"
	"playbook/responses"
	"playbook/services"
	"playbook/utils"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PickupLineListResponse struct {
	Total       int                                       `json:"total"`
	Page        int                                       `json:"page"`
	PickupLines []*responses.SinglePickupLineInfoResponse `json:"pickup_lines"`
}

type PickupLineController struct {
	pickupLineService *services.PickupLineService
}

// CreatePickupLine godoc
// @Summary      Create a PickupLine
// @Description  Create a PickupLine
// @Tags         Pickup-Line
// @Accept       json
// @Produce      json
//
// @Param request body requests.PickupLineBodyRequest true "PickupLine to add"
//
// @Security Bearer
// @Success      200  {object}  responses.SinglePickupLineInfoResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /pickup-lines [post]
func (u *PickupLineController) CreatePickupLine(ctx *gin.Context) {
	pickupLine := &requests.PickupLineBodyRequest{}

	userUUID, ok := utils.ExtractUserUuidFromRequest(ctx)
	if !ok {
		return
	}

	claims := jwt.ExtractClaims(ctx)
	username := claims[constants.JWT_USERNAME_KEY].(string)
	if username == "" {
		ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(500, "Missing username"))
		return
	}

	displayName := claims[constants.JWT_DISPLAY_NAME_KEY].(string)
	if displayName == "" {
		ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(400, "Missing display name"))
		return
	}

	user := &entities.User{
		ID:          userUUID,
		Username:    username,
		DisplayName: displayName,
	}
	err := ctx.ShouldBind(pickupLine)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(500, "error while parsing the pickupLine"))
		return
	}

	createdPickupLine, err := u.pickupLineService.Create(pickupLine, user)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(500, "error while creating the pickupLine"))
		return
	}

	ctx.JSON(http.StatusOK, responses.NewSinglePickupLineInfoResponse(createdPickupLine))
}

// DeletePickupLine godoc
// @Summary      Delete a PickupLine
// @Description  Delete a PickupLine
// @Tags         Pickup-Line
// @Accept       json
// @Produce      json
//
// @Param        pickupLineId   path      string  true  "PickupLine ID"
//
// @Security Bearer
// @Success      200  {object}  responses.GenericSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /pickup-lines/{pickupLineId} [delete]
func (u *PickupLineController) DeletePickupLine(ctx *gin.Context) {
	userUUID, ok := utils.ExtractUserUuidFromRequest(ctx)
	if !ok {
		return
	}

	pickupLineId := ctx.Param(constants.PickupLineId)
	if pickupLineId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "Missing pickupLineId"})
		return
	}

	pickupLineUUID, err := uuid.Parse(pickupLineId)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "invalid pickupLineId"})
		return
	}

	err = u.pickupLineService.Delete(pickupLineUUID, userUUID)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Could not delete pickup line"})
		return
	}

	ctx.JSON(http.StatusOK, responses.NewGenericSuccessResponse())
}

// UpdatePickupLine godoc
// @Summary      Update a PickupLine
// @Description  Update a PickupLine
// @Tags         Pickup-Line
// @Accept       json
// @Produce      json
//
// @Param        pickupLineId   path      string  true  "PickupLine ID"
// @Param request body requests.PickupLineBodyRequest true "PickupLine to update"
//
// @Security Bearer
// @Success      200  {object}  responses.SinglePickupLineInfoResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /pickup-lines/{pickupLineId} [put]
func (u *PickupLineController) UpdatePickupLine(ctx *gin.Context) {
	pickupLine := &requests.PickupLineBodyRequest{}

	pickupLineId := ctx.Param(constants.PickupLineId)
	if pickupLineId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "Missing pickupLineId"})
		return
	}

	userUUID, ok := utils.ExtractUserUuidFromRequest(ctx)
	if !ok {
		return
	}

	err := ctx.ShouldBind(pickupLine)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(500, "error while parsing the pickupLine"))
		return
	}

	pickupLineUuid, err := uuid.Parse(pickupLineId)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "invalid pickupLineId"})
		return
	}

	updatedPickupLine, err := u.pickupLineService.Update(pickupLineUuid, pickupLine, userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "could not find pickupLine with given id"})
			return
		}

		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Could not update pickupline"})
		return
	}

	ctx.JSON(http.StatusOK, responses.NewSinglePickupLineInfoResponse(updatedPickupLine))
}

// UpdatePickupLineReaction godoc
// @Summary      Update a PickupLine Reaction
// @Description  Update a PickupLine Reaction
// @Tags         Pickup-Line
// @Accept       json
// @Produce      json
//
// @Param        pickupLineId   path      string  true  "PickupLine ID"
// @Param request body entities.Reaction true "Reaction to the PickupLine"
//
// @Security Bearer
// @Success      200  {object}  responses.GenericSuccessResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Router       /pickup-lines/{pickupLineId}/reaction [put]
// @Failure      500  {object}  responses.ErrorResponse
func (u *PickupLineController) UpdatePickupLineReaction(ctx *gin.Context) {
	pickupLineId := ctx.Param(constants.PickupLineId)
	if pickupLineId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Missing pickupLineId"})
		return
	}

	userUUID, ok := utils.ExtractUserUuidFromRequest(ctx)
	if !ok {
		return
	}

	pickupLineUUID, err := uuid.Parse(pickupLineId)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "invalid pickupLineId"})
		return
	}

	if !u.pickupLineService.CanUserSeePickupLine(pickupLineUUID, userUUID) {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "Pickupline is not visible"})
		return
	}

	reaction := &entities.Reaction{}
	err = ctx.ShouldBind(reaction)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(400, "error while parsing the pickupLine"))
		return
	}
	reaction.UserId = userUUID
	reaction.PickupLineId = pickupLineUUID

	err = u.pickupLineService.UpdateReactionByUser(userUUID, pickupLineUUID, reaction)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "could not update Pickupline reaction"})
		return
	}

	ctx.JSON(http.StatusOK, responses.NewGenericSuccessResponse())
}

// GetPickupLineById godoc
// @Summary      Get a PickupLine by id
// @Description  Get a PickupLine by id
// @Tags         Pickup-Line
// @Accept       json
// @Produce      json
//
// @Param        pickupLineId   path      string  true  "PickupLine ID"
//
// @Security Bearer
// @Success      200  {object}  responses.SinglePickupLineInfoResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /pickup-lines/{pickupLineId} [get]
func (t *PickupLineController) GetPickupLineById(ctx *gin.Context) {
	pickupLineId := ctx.Param(constants.PickupLineId)
	if pickupLineId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "Missing pickupLineId"})
		return
	}

	pickupLineUuid, err := uuid.Parse(pickupLineId)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "invalid pickupLineId"})
		return
	}

	userUUID, ok := utils.ExtractUserUuidFromRequest(ctx)
	if !ok {
		return
	}

	pickupLine, err := t.pickupLineService.GetByIdAndUserId(pickupLineUuid, userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "could not find pickupLine with given id"})
			return
		}

		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "error while retrieving pickupLine"})
		return
	}

	statistics, err := t.pickupLineService.GetStatisticsById(pickupLineUuid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "error while retrieving pickupLine"})
		return
	}
	pickupLine.Statistics = statistics

	if len(pickupLine.Reactions) > 0 {
		pickupLine.UserReaction = pickupLine.Reactions[0]
	} else {
		pickupLine.UserReaction = &entities.Reaction{
			PickupLineId: pickupLine.ID,
			UserId:       userUUID,
			Starred:      false,
			Vote:         entities.None,
		}
	}

	ctx.JSON(http.StatusOK, responses.NewSinglePickupLineInfoResponse(pickupLine))
}

// GetPickupLineList godoc
// @Summary      Get a PickupLine List
// @Description  Get a PickupLine List
// @Tags         Pickup-Line
// @Accept       json
// @Produce      json
//
// @Param        filters query  requests.PickupLineFilters  false  "Filters"
//
// @Security Bearer
// @Success      200  {object}  PickupLineListResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /pickup-lines [get]
func (t *PickupLineController) GetPickupLineList(ctx *gin.Context) {
	userUUID, ok := utils.ExtractUserUuidFromRequest(ctx)
	if !ok {
		return
	}

	filters := &requests.PickupLineFilters{}
	ctx.BindQuery(&filters)

	searchFilters := getElasticFilterFromRequestFilter(filters)
	elasticResponse, err := t.pickupLineService.GetList(userUUID, searchFilters)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "Could not retrieve pickupLine list"})
		return
	}

	response := &PickupLineListResponse{
		Total:       elasticResponse.Total,
		Page:        filters.Page,
		PickupLines: elasticResponse.Users,
	}

	ctx.JSON(http.StatusOK, response)
}

func (t *PickupLineController) GetPickupLineFeed(ctx *gin.Context) {
	userUUID, ok := utils.ExtractUserUuidFromRequest(ctx)
	if !ok {
		return
	}

	filters := &requests.PickupLineFilters{}
	ctx.BindQuery(&filters)

	searchFilters := getElasticFilterFromRequestFilter(filters)
	elasticResponse, err := t.pickupLineService.GetFeed(userUUID, searchFilters)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "Could not Get PickupLine feed"})
		return
	}

	response := &PickupLineListResponse{
		Total:       elasticResponse.Total,
		Page:        filters.Page,
		PickupLines: elasticResponse.Users,
	}

	ctx.JSON(http.StatusOK, response)
}

func getElasticFilterFromRequestFilter(rawFilters *requests.PickupLineFilters) *filters.PickupLineQueryFilters {
	userId, err := uuid.Parse(rawFilters.UserId)
	if err != nil {
		userId = uuid.Nil
	}

	return &filters.PickupLineQueryFilters{
		Page:              rawFilters.Page,
		Title:             rawFilters.Title,
		Starred:           rawFilters.Starred,
		Tags:              rawFilters.Tags,
		Visible:           rawFilters.Visible,
		OnlyUpvoted:       rawFilters.OnlyUpvoted,
		Content:           rawFilters.Content,
		SuccessPercentage: rawFilters.SuccessPercentage,
		UserId:            userId,
		SortingType:       rawFilters.SortingType,
	}
}

func CreatePickupLineController() *PickupLineController {
	db := database.GetDatabaseConnection()

	tagRepository := repositories.CreateTagRepository(db)
	pickupLineRepository := repositories.CreatePickupLineRepository(db)
	elasticSearchMapper := services.GetNewElasticSearchMapper(tagRepository, pickupLineRepository)
	elasticSearchWrapper := services.GetElasticSearchWrapper(elasticSearchMapper)
	mapper := mappers.GetNewDtoMapper()
	tagService := services.CreateTagService(tagRepository, elasticSearchWrapper, mapper)
	pickupLineService := services.CreatePickupLineService(pickupLineRepository, tagService, elasticSearchWrapper, mapper)

	return &PickupLineController{pickupLineService: pickupLineService}
}
