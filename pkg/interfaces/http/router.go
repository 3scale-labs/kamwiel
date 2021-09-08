package http

import (
	"fmt"
	"github.com/3scale-labs/kamwiel/pkg/adapters/kuadrant"
	"github.com/3scale-labs/kamwiel/pkg/interfaces/http/handlers"
	"github.com/3scale-labs/kamwiel/pkg/services/api"
	"github.com/gin-gonic/gin"
	"os"
)

var router = gin.Default()

func urlMappings() {
	apiHandler := handlers.NewAPIHandler(
		api.NewService(
			kuadrant.NewKuadrantRepository(kuadrant.Client)))

	router.GET("/ping", handlers.Ping)
	router.GET("/apis", apiHandler.List)
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
