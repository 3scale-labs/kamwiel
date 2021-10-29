package handlers

import (
	"fmt"
	"github.com/3scale-labs/kamwiel/pkg/services/api"
	"github.com/gin-gonic/gin"
	apiErrors "k8s.io/apimachinery/pkg/api/errors"
	"net/http"
)

type APIHandler interface {
	Get(*gin.Context)
	List(*gin.Context)
	GetListState(*gin.Context)
	UpdateListState(*gin.Context)
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
	apiObj, getErr := h.service.GetAPI(ctx, name)

	if getErr != nil && apiErrors.IsNotFound(getErr) {
		fmt.Println("API not found", getErr)
		ctx.JSON(http.StatusNotFound, fmt.Sprintf("API %s not found", name))
		return
	} else if getErr != nil {
		ctx.JSON(http.StatusInternalServerError, getErr)
		return
	}

	ctx.JSON(http.StatusOK, apiObj)
}

func (h *apiHandler) List(ctx *gin.Context) {
	list, getErr := h.service.ListAPI(ctx)
	if getErr != nil && apiErrors.IsNotFound(getErr) {
		fmt.Println("No APIs were found", getErr)
		ctx.JSON(http.StatusNotFound, "No APIs were found")
		return
	} else if getErr != nil {
		ctx.JSON(http.StatusInternalServerError, getErr)
		return
	}

	ctx.JSON(http.StatusOK, list)
}

func (h *apiHandler) GetListState(ctx *gin.Context) {
	apiListState, err := h.service.GetAPIListState(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apiListState)
}

func (h *apiHandler) UpdateListState(ctx *gin.Context) {
	hash := ctx.Param("hash")
	if len(hash) == 0 { // this check might be a candidate to extract and reuse
		ctx.JSON(http.StatusBadRequest, "Missing param `hash`")
		return
	}
	err := h.service.UpdateAPIListState(ctx, hash)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusNoContent, "API List State updated successfully")
}
