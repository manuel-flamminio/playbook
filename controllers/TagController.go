package controllers

import (
	"errors"
	"log"
	"net/http"
	"playbook/constants"
	"playbook/database"
	"playbook/entities"
	"playbook/mappers"
	"playbook/repositories"
	"playbook/requests"
	"playbook/responses"
	"playbook/services"
	"playbook/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagController struct {
	tagService services.TagServiceInterface
}

// CreateTag godoc
// @Summary      Create a Tag
// @Description  Create a Tag
// @Tags         Tag
// @Accept       json
// @Produce      json
//
// @Param request body requests.TagBodyRequest true "tag to add"
//
// @Security Bearer
// @Success      200  {object}  entities.Tag
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /tags [post]
func (u *TagController) CreateTag(ctx *gin.Context) {
	tag := &requests.TagBodyRequest{}

	userUUID, ok := utils.ExtractUserUuidFromRequest(ctx)
	if !ok {
		return
	}

	err := ctx.ShouldBind(tag)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(400, "Error while parsing the Tag"))
		return
	}

	createdTag, err := u.tagService.Create(tag, userUUID)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(500, "Error while saving the Tag"))
		return
	}

	ctx.JSON(http.StatusOK, createdTag)
}

// DeleteTag godoc
// @Summary      Delete a Tag
// @Description  Delete a Tag
// @Tags         Tag
// @Accept       json
// @Produce      json
//
// @Param        tagId   path      string  true  "Tag ID"
//
// @Security Bearer
// @Success      200  {object}  responses.GenericSuccessResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /tags/{tagId} [delete]
func (u *TagController) DeleteTag(ctx *gin.Context) {
	tag := &entities.Tag{}

	userUUID, ok := utils.ExtractUserUuidFromRequest(ctx)
	if !ok {
		return
	}

	tag.UserId = userUUID

	tagId := ctx.Param(constants.TagId)
	if tagId == "" {
		ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(400, "invalid Tag Id"))
		return
	}

	tagUUID, err := uuid.Parse(tagId)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(400, "invalid Tag Id"))
		return
	}

	tag.ID = tagUUID
	err = u.tagService.Delete(tag)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(500, "Error while deleting the Tag"))
		return
	}

	ctx.JSON(http.StatusOK, responses.NewGenericSuccessResponse())
}

// UpdateTag godoc
// @Summary      Update a Tag
// @Description  update a Tag (the update is in state mode, so you need to pass the entire object)
// @Tags         Tag
// @Accept       json
// @Produce      json
//
// @Param        tagId   path      string  true  "Tag ID"
// @Param request body requests.TagBodyRequest true "tag to update"
//
// @Security Bearer
// @Success      200  {object}  entities.Tag
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      404  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /tags/{tagId} [put]
func (u *TagController) UpdateTag(ctx *gin.Context) {
	tag := &requests.TagBodyRequest{}

	tagId := ctx.Param(constants.TagId)
	if tagId == "" {
		ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(400, "invalid Tag Id"))
		return
	}

	err := ctx.ShouldBind(tag)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(400, "Error while parsing the Tag"))
		return
	}

	updatedTag, err := u.tagService.Update(tagId, tag) //todo: update tag not from user
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, responses.NewErrorResponse(404, "Could not find Tag with given id"))
			return
		}

		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(500, "Error while deleting the Tag"))
		return
	}

	ctx.JSON(http.StatusOK, updatedTag)
}

// GetTagList godoc
// @Summary      Get User Tag List
// @Description  Get the logged user tag list
// @Tags         Tag
// @Produce      json
//
// @Security Bearer
// @Success      200  {array}  entities.Tag
// @Failure      401  {object} responses.ErrorResponse
// @Failure      500  {object} responses.ErrorResponse
// @Router       /tags [get]
func (t *TagController) GetTagList(ctx *gin.Context) {
	userUUID, ok := utils.ExtractUserUuidFromRequest(ctx)
	if !ok {
		return
	}

	tags, err := t.tagService.GetListByUserId(userUUID)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(500, "Error while retrieving the Tag list"))
		return
	}

	ctx.JSON(http.StatusOK, tags)
}

func CreateTagController() *TagController {
	db := database.GetDatabaseConnection()

	repository := repositories.CreateTagRepository(db)
	pickupLineRepository := repositories.CreatePickupLineRepository(db)
	elasticSearchMapper := services.GetNewElasticSearchMapper(repository, pickupLineRepository)
	elasticSearchWrapper := services.GetElasticSearchWrapper(elasticSearchMapper)
	mapper := mappers.GetNewDtoMapper()
	service := services.CreateTagService(repository, elasticSearchWrapper, mapper)
	return &TagController{tagService: service}
}
