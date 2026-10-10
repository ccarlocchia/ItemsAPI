package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"project/cmd/api/handlers"
	"project/internal/model"
	"project/internal/repositories"
	"project/internal/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}
func setupIntegrationRouter(t *testing.T) *gin.Engine {
	t.Helper()

	// Ensure mocked_items.json exists
	if _, err := os.Stat("mocked_items.json"); err != nil {
		t.Skip("mocked_items.json not found, skipping integration test")
	}

	// Initialize real components
	repo := repositories.NewRepository()
	service := services.NewItemsService(&repo)
	handler := handlers.NewHandler(service)

	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler.SetupRoutes(router)

	return router
}

func TestIntegration_GetItemsByIDs(t *testing.T) {
	router := setupIntegrationRouter(t)

	t.Run("valid ID returns 200 and correct item", func(t *testing.T) {
		body := handlers.RequestBody{
			ItemID: []handlers.ItemData{{ID: "123abc"}},
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/api/items", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d, body=%s", rec.Code, rec.Body.String())
		}

		var resp struct {
			Data []model.ItemDTO `json:"data"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(resp.Data) != 1 {
			t.Fatalf("expected 1 item, got %d", len(resp.Data))
		}
		item := resp.Data[0]
		if item.ID != "123abc" || item.Name != "Smartphone X" {
			t.Errorf("expected item Smartphone X with ID 123abc, got %+v", item)
		}
	})

	t.Run("non-existent ID returns 404", func(t *testing.T) {
		body := handlers.RequestBody{
			ItemID: []handlers.ItemData{{ID: "nonexistent"}},
		}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/api/items", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected status 404, got %d, body=%s", rec.Code, rec.Body.String())
		}
		if !bytes.Contains(rec.Body.Bytes(), []byte("no items found")) {
			t.Errorf("expected error message 'no items found', got %s", rec.Body.String())
		}
	})

	t.Run("empty body returns 400", func(t *testing.T) {
		body := handlers.RequestBody{ItemID: []handlers.ItemData{}}
		bodyBytes, _ := json.Marshal(body)

		req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/api/items", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected status 400, got %d, body=%s", rec.Code, rec.Body.String())
		}
		if !bytes.Contains(rec.Body.Bytes(), []byte("invalid request body")) {
			t.Errorf("expected error message 'invalid request body', got %s", rec.Body.String())
		}
	})
}
