package handlers

import (
	"context"
	"net/http"
	"project/internal/model"
	"strings"

	"github.com/gin-gonic/gin"
)

type ItemsService interface {
	GetItemsByIDs(ctx context.Context, mappedItems []string) []model.ItemDTO
}

type Handler struct {
	service ItemsService
}

func NewHandler(service ItemsService) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetData(ctx *gin.Context) {
	var body RequestBody
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	ids := extractItemIDs(body)
	if len(ids) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	result := h.service.GetItemsByIDs(ctx, ids)
	if len(result) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "no items found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": result})
}

type RequestBody struct {
	ItemID []ItemData `json:"data" binding:"required"`
}

type ItemData struct {
	ID string `json:"id"`
}

func (h *Handler) SetupRoutes(router *gin.Engine) {
	router.POST("/v1/api/items", h.GetData)
}

func extractItemIDs(body RequestBody) []string {
	ids := make([]string, 0, len(body.ItemID))
	for _, item := range body.ItemID {
		id := strings.TrimSpace(item.ID)
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}
