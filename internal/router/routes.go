package router

import (
	"automation-hub-backend/docs"
	"automation-hub-backend/internal/automation"
	"automation-hub-backend/internal/config"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initializeRoutes(router *gin.Engine) error {
	relativePathV1 := config.AppConfig.BaseUrl + "/v1"
	docs.SwaggerInfo.BasePath = relativePathV1
	v1 := router.Group(relativePathV1)
	{
		autoHandler := automation.DefaultHandler()
		err := initializeAutomationsRoutes(v1, autoHandler)
		if err != nil {
			return err
		}
		v1.GET("/images/:imageName", autoHandler.ImageHandler)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	return nil
}

func initializeAutomationsRoutes(apiVersion *gin.RouterGroup, autoHandler *automation.Handler) error {
	automations := apiVersion.Group("/automations")
	{
		automations.GET("/", autoHandler.GetAll)
		automations.GET("/:id", autoHandler.GetByID)
		automations.POST("/", autoHandler.Create)
		automations.PATCH("/", autoHandler.Update)
		automations.DELETE("/:id", autoHandler.DeleteByID)
		automations.GET("/:id1/swap/:id2", autoHandler.SwapPosition)
	}

	return nil
}
