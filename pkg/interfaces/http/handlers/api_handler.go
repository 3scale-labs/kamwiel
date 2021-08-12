package handlers

import (
	"fmt"
	"github.com/3scale-labs/kamwiel/pkg/domain/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type APIHandler interface {
	Get(*gin.Context)
}

type apiHandler struct {
	service api.Service
}

func NewAPIHandler(service api.Service) APIHandler {
	return &apiHandler{
		service: service,
	}
}

func (h *apiHandler) Get(ctx *gin.Context) {
	name := ctx.Param("name")
	if len(name) == 0 {
		ctx.JSON(http.StatusBadRequest, "Missing param `name`")
		return
	}
	api, getErr := h.service.GetAPI(name)
	if getErr != nil {
		fmt.Println("API not found", getErr)
		ctx.JSON(http.StatusNotFound, "API not found")
		return
	}

	ctx.JSON(http.StatusOK, api)
}
