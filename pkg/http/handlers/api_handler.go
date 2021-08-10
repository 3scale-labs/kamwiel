package handlers

import (
	"fmt"
	"github.com/3scale-labs/kamwiel/pkg/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Get(ctx *gin.Context) {
	name := ctx.Param("name")
	if len(name) == 0 {
		ctx.JSON(http.StatusBadRequest, "Missing param `name`")
		return
	}
	api, getErr := services.APIService.GetAPI(name)
	if getErr != nil {
		fmt.Println("API not found", getErr)
		ctx.JSON(http.StatusNotFound, "API not found")
		return
	}

	ctx.JSON(http.StatusOK, api)
}
