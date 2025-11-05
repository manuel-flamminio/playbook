package middlewares

import (
	"os"
	"playbook/constants"
	"playbook/entities"
	"playbook/services"

	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	userService *services.UserService
}

func Authorizator(ctx *gin.Context) {

}

func (a *AuthMiddleware) Authenticator(ctx *gin.Context) (interface{}, error) {
	user := &entities.User{}
	err := ctx.ShouldBind(user)
	if err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	user, err = a.userService.Login(user)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return user, nil
}

func (a *AuthMiddleware) GetMiddleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(
		&jwt.GinJWTMiddleware{
			Realm:         "Playbook",
			Key:           []byte(os.Getenv(constants.JWT_SECRET_KEY)),
			Timeout:       time.Hour,
			MaxRefresh:    time.Hour,
			PayloadFunc:   payloadFunc(),
			Authenticator: a.Authenticator,
			TokenLookup:   "header: Authorization, query: token",
			TokenHeadName: "Bearer",
			TimeFunc:      time.Now,
		},
	)
}

func payloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if u, ok := data.(*entities.User); ok {
			return jwt.MapClaims{
				jwt.IdentityKey:                u.ID,
				constants.JWT_USERNAME_KEY:     u.Username,
				constants.JWT_DISPLAY_NAME_KEY: u.DisplayName,
			}
		}
		return jwt.MapClaims{}
	}
}

func HandlerMiddleware(authMiddleware *jwt.GinJWTMiddleware) gin.HandlerFunc {
	return func(context *gin.Context) {
		errInit := authMiddleware.MiddlewareInit()
		if errInit != nil {
			panic("could not handle auth middleware")
		}
	}
}

func CreateAuthMiddleware(userService *services.UserService) *AuthMiddleware {
	return &AuthMiddleware{userService: userService}
}
