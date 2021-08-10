package http

import (
	"fmt"
	"github.com/3scale-labs/kamwiel/pkg/http/handlers"
	"github.com/gin-gonic/gin"
	"os"
)

var router = gin.Default()

func urlMappings() {
	router.GET("/ping", handlers.Ping)
	router.GET("/apis/:name", handlers.Get)
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
