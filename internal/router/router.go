package router

import (
	"automation-hub-backend/internal/config"
	"github.com/gin-gonic/gin"
)

func Initialize() error {
	// initialize Router
	router := gin.Default()

	// initialize routes
	err := initializeRoutes(router)
	if err != nil {
		return err
	}

	// run server
	port := ":" + config.Config.ServerPort
	err = router.Run(port)
	if err != nil {
		return err
	}

	return nil
}
