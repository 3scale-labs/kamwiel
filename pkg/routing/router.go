package routing

import (
	"fmt"
	"github.com/3scale-labs/kamwiel/pkg/controllers/api"
	"github.com/3scale-labs/kamwiel/pkg/controllers/ping"
	"github.com/gin-gonic/gin"
	"os"
)

var router = gin.Default()

func urlMappings() {
	router.GET("/ping", ping.Ping)
	router.GET("/apis/:name", api.Get)
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
