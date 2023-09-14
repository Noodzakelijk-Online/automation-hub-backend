package router

import (
	"automation-hub-backend/docs"
	"automation-hub-backend/internal/config"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initializeRoutes(router *gin.Engine) error {
	relativePathV1 := config.Config.BaseUrl + "/v1"
	docs.SwaggerInfo.BasePath = relativePathV1
	v1 := router.Group(relativePathV1)
	{
		// initialize auth routes
		err := initializeAutomationsRoutes(v1)
		if err != nil {
			return err
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return nil
}

func initializeAutomationsRoutes(apiVersion *gin.RouterGroup) error {
	automations := apiVersion.Group("/automations")
	{
		// initialize automations routes
		//autoHandler := authentication.NewHandler(autoService)
		// fake handler function
		autoHandler := func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		}

		automations.GET("/", autoHandler)
		automations.POST("/", autoHandler)
		automations.PATCH("/:id", autoHandler)
		automations.GET("/:id", autoHandler)
		automations.DELETE("/:id", autoHandler)
		automations.POST("/:id1/swap/:id2", autoHandler)
	}

	return nil
}
