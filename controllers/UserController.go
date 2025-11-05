package controllers

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"playbook/constants"
	"playbook/database"
	"playbook/filters"
	"playbook/mappers"
	"playbook/repositories"
	"playbook/requests"
	"playbook/responses"
	"playbook/services"
	"playbook/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController struct {
	userService       *services.UserService
	pickupLineService *services.PickupLineService
}

// CreateUser godoc
// @Summary      Create a User
// @Description  Create a User
// @Tags         User
// @Accept       json
// @Produce      json
//
// @Param request body requests.CreateUserRequest true "user to add"
//
// @Success      200  {object}  responses.BaseUserInfoResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /user [post]
func (u *UserController) CreateUser(ctx *gin.Context) {
	user := &requests.CreateUserRequest{}
	err := ctx.ShouldBind(user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, responses.NewErrorResponse(400, "Error while parsing the User"))
		return
	}

	if user.EncodedUserImage != "" {
		decoded, err := base64.StdEncoding.DecodeString(user.EncodedUserImage)
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(500, "Error while parsing the image"))
			return
		}

		image, err := jpeg.Decode(bytes.NewReader(decoded))
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(500, "Error while parsing the image"))
			return
		}

		f, err := os.OpenFile(os.Getenv(constants.STORAGE_PATH)+"/"+user.Username+".jpeg", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(500, "Error while parsing the image"))
			return
		}
		err = jpeg.Encode(f, image, &jpeg.Options{Quality: 75})
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(500, "Error while parsing the image"))
			return
		}
		f.Close()
		user.EncodedUserImage = ""
	}

	createdUser, err := u.userService.Create(user)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, responses.NewErrorResponse(500, "Error while creating the user"))
		return
	}

	ctx.JSON(http.StatusOK, responses.NewBaseUserInfoResponse(createdUser))
}

// UpdateUserDisplayName godoc
// @Summary      Update the User display name
// @Description  Update the User display name
// @Tags         User
// @Accept       json
// @Produce      json
//
// @Param request body requests.UpdateUserDisplayNameRequest true "New display name"
//
// @Security Bearer
// @Success      200  {object}  responses.BaseUserInfoResponse
// @Failure      400  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /user/name [put]
func (u *UserController) UpdateUserDisplayName(ctx *gin.Context) { //todo: update by query pickuplines display name
	user := &requests.UpdateUserDisplayNameRequest{}

	userId, ok := utils.ExtractUserUuidFromRequest(ctx)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "Missing UserId"})
		return
	}

	err := ctx.ShouldBind(user)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "error while parsing the user"})
		return
	}

	updatedUser, err := u.userService.UpdateDisplayName(userId, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "could not find user with given id "})
			return
		}

		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "error while creating the user"})
		return
	}
	updatedUser.Password = ""

	ctx.JSON(http.StatusOK, updatedUser)
}

// GetUsernameInfo godoc
// @Summary      Get User Info
// @Description  Get User Info
// @Tags         User
// @Accept       json
// @Produce      json
//
// @Security Bearer
// @Success      200  {object}  responses.BaseUserInfoResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /user/info [get]
func (u *UserController) GetUsernameInfo(ctx *gin.Context) { //todo: add user general statistics
	user, ok := utils.ExtractUserFromRequest(ctx)
	if !ok {
		return
	}

	ctx.JSON(http.StatusOK, responses.NewBaseUserInfoResponse(user))
}

// DeleteUser godoc
// @Summary      Delete a user and all related data
// @Description  Delete a user and all related data
// @Tags         User
// @Accept       json
// @Produce      json
//
// @Security Bearer
// @Success      200  {object}  responses.GenericSuccessResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /user [delete]
func (u *UserController) DeleteUser(ctx *gin.Context) {
	user, ok := utils.ExtractUserFromRequest(ctx)
	if !ok {
		return
	}

	if err := u.pickupLineService.DeleteByUser(user.ID); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, &responses.ErrorResponse{Code: 500, Message: "error while deleting user data, contact the administrator"})
		return
	}

	if err := u.userService.Delete(user); err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, &responses.ErrorResponse{Code: 500, Message: "error while deleting the user"})
		return
	}

	ctx.JSON(http.StatusOK, responses.NewGenericSuccessResponse())
}

// SearchUsers godoc
// @Summary      Search users
// @Description  Search users
// @Tags         User
// @Accept       json
// @Produce      json
//
// @Param        filters query  requests.UserFilters  false  "Filters"
//
// @Security Bearer
// @Success      200  {object}  responses.ElasticSearchUserListResponse
// @Failure      401  {object}  responses.ErrorResponse
// @Failure      500  {object}  responses.ErrorResponse
// @Router       /user [get]
func (t *UserController) GetSearchUser(ctx *gin.Context) {
	filters := &requests.UserFilters{}
	ctx.BindQuery(&filters)

	searchFilters := getUserElasticFilterFromRequestFilter(filters)
	elasticResponse, err := t.userService.SearchUser(searchFilters)
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Could not retrieve users"})
		return
	}
	elasticResponse.Page = filters.Page

	ctx.JSON(http.StatusOK, elasticResponse)
}

func (c *UserController) GetUserService() *services.UserService {
	return c.userService
}

func getUserElasticFilterFromRequestFilter(rawFilters *requests.UserFilters) *filters.UserQueryFilters {
	return &filters.UserQueryFilters{
		Page:        rawFilters.Page,
		Username:    rawFilters.Username,
		DisplayName: rawFilters.DisplayName,
	}
}

func CreateUserController() *UserController {
	db := database.GetDatabaseConnection()

	tagRepository := repositories.CreateTagRepository(db)
	pickupLineRepository := repositories.CreatePickupLineRepository(db)
	elasticSearchMapper := services.GetNewElasticSearchMapper(tagRepository, pickupLineRepository)
	elasticSearchWrapper := services.GetElasticSearchWrapper(elasticSearchMapper)
	mapper := mappers.GetNewDtoMapper()
	tagService := services.CreateTagService(tagRepository, elasticSearchWrapper, mapper)
	pickupLineService := services.CreatePickupLineService(pickupLineRepository, tagService, elasticSearchWrapper, mapper)

	repository := repositories.CreateUserRepository(db)
	service := services.CreateUserService(repository, mapper, elasticSearchWrapper)
	return &UserController{userService: service, pickupLineService: pickupLineService}
}
