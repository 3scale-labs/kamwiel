package routing

import (
	"fmt"
	"github.com/3scale-labs/kamwiel/controllers/ping"
	"github.com/gin-gonic/gin"
)

var router = gin.Default()

func urlMappings() {
	router.GET("/ping", ping.Ping)
}

func Start() {
	urlMappings()
	fmt.Println("Kamwiel starting...")
	router.Run(":8080")
}
