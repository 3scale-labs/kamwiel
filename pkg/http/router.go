package http

import (
	"fmt"
	"github.com/3scale-labs/kamwiel/pkg/http/handlers"
	"github.com/3scale-labs/kamwiel/pkg/repositories"
	"github.com/3scale-labs/kamwiel/pkg/services"
	"github.com/gin-gonic/gin"
	"os"
)

var router = gin.Default()

func urlMappings() {
	apiHandler := handlers.NewAPIHandler(
		services.NewAPIService(
			repositories.NewKuadrantRepository()))

	router.GET("/ping", handlers.Ping)
	router.GET("/apis/:name", apiHandler.Get)
}

func Start() {
	urlMappings()
	port, ok := os.LookupEnv("PORT")
	if !ok {
		panic("ENV PORT is not present")
	} else {
		fmt.Println("Kamwiel listening on port " + port)
		router.Run(":" + port)
	}
}
