package utils

import (
	"net/http"
	"playbook/constants"
	"playbook/entities"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	jwt "github.com/appleboy/gin-jwt/v2"
)

func ExtractUserUuidFromRequest(ctx *gin.Context) (uuid.UUID, bool) {
	claims := jwt.ExtractClaims(ctx)
	userId := claims[jwt.IdentityKey].(string)
	if userId == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "Missing userId"})
		return uuid.Nil, false
	}

	userUUID, err := uuid.Parse(userId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "invalid userId"})
		return uuid.Nil, false
	}

	return userUUID, true
}

func ExtractUserFromRequest(ctx *gin.Context) (*entities.User, bool) {
	userUUID, ok := ExtractUserUuidFromRequest(ctx)
	if !ok {
		return nil, false
	}

	claims := jwt.ExtractClaims(ctx)
	username := claims[constants.JWT_USERNAME_KEY].(string)
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "Missing username"})
		return nil, false
	}

	displayName := claims[constants.JWT_DISPLAY_NAME_KEY].(string)
	if displayName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 500, "message": "Missing display name"})
		return nil, false
	}

	return &entities.User{
		ID:          userUUID,
		Username:    username,
		DisplayName: displayName,
	}, true
}
