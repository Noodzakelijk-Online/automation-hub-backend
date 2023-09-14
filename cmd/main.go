package main

import (
	"automation-hub-backend/internal/config"
	"automation-hub-backend/internal/router"
)

func main() {
	config.Init()

	err := router.Initialize()
	if err != nil {
		panic(err)
	}

}
