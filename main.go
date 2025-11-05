package main

import (
	"fmt"
	"os"
	"playbook/constants"
	"playbook/controllers"
	"playbook/database"
	_ "playbook/docs"
	"playbook/middlewares"
	"playbook/seeder"
	"playbook/services"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title The playbook
// @version 1.0
// @description     The playbook
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      playbook.example.com
// @BasePath  /api

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @authorizationurl https://playbook.example.com/login

// @schemes https http

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	migrateAll()
	r := gin.Default()

	userController := controllers.CreateUserController()
	tagController := controllers.CreateTagController()
	pickupLineController := controllers.CreatePickupLineController()
	authMiddleware := middlewares.CreateAuthMiddleware(userController.GetUserService())

	middleware, err := authMiddleware.GetMiddleware()
	if err != nil {
		panic("could not instantiate authentication middleware")
	}

	r.Static("/api/images", os.Getenv(constants.STORAGE_PATH))

	r.POST("/api/login", middleware.LoginHandler)
	r.POST("/api/user", userController.CreateUser)
	r.POST("/api/seed/:secret", seeder.SeedData)

	r.Use(middlewares.HandlerMiddleware(middleware))

	api := r.Group("/api", middleware.MiddlewareFunc())
	api.PUT("/user/name", userController.UpdateUserDisplayName)
	api.GET("/user/info", userController.GetUsernameInfo)
	api.GET("/user", userController.GetSearchUser)
	api.DELETE("/user", userController.DeleteUser) //todo clean session after user deletion

	api.POST("/tags", tagController.CreateTag)
	api.PUT(fmt.Sprintf("/tags/:%s", constants.TagId), tagController.UpdateTag)
	api.DELETE(fmt.Sprintf("/tags/:%s", constants.TagId), tagController.DeleteTag)
	api.GET("/tags", tagController.GetTagList)

	api.POST("/pickup-lines", pickupLineController.CreatePickupLine)
	api.PUT(fmt.Sprintf("/pickup-lines/:%s", constants.PickupLineId), pickupLineController.UpdatePickupLine)
	api.PUT(fmt.Sprintf("/pickup-lines/:%s/reaction", constants.PickupLineId), pickupLineController.UpdatePickupLineReaction)
	api.DELETE(fmt.Sprintf("/pickup-lines/:%s", constants.PickupLineId), pickupLineController.DeletePickupLine)
	api.GET(fmt.Sprintf("/pickup-lines/:%s", constants.PickupLineId), pickupLineController.GetPickupLineById)
	api.GET("/pickup-lines", pickupLineController.GetPickupLineList)
	api.GET("/pickup-lines/feed", pickupLineController.GetPickupLineFeed) //todo: add algorithm

	auth := r.Group("/auth", middleware.MiddlewareFunc())
	auth.GET("/refresh_token", middleware.RefreshHandler)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Run()
}

func migrateAll() {
	db := database.ConnectToDatabase()
	database.MigrateDatabase(db)

	e := services.GetElasticSearchWrapper(nil)
	exists, err := e.ExistsIndex(e.PickupLineIndexName)
	if err != nil {
		panic("could not check pickupLine index existance")
	}
	if !exists {
		if err := e.CreatePickupLineIndex(); err != nil {
			panic("could not create pickupLine index")
		}
	}

	exists, err = e.ExistsIndex(e.TagsIndexName)
	if err != nil {
		panic("could not check tags index existance")
	}
	if !exists {
		if err := e.CreateTagIndex(); err != nil {
			panic("could not create tags index")
		}
	}

	exists, err = e.ExistsIndex(e.UserIndexName)
	if err != nil {
		panic("could not check user index existance")
	}
	if !exists {
		if err := e.CreateUserIndex(); err != nil {
			panic("could not create user index")
		}
	}
}
