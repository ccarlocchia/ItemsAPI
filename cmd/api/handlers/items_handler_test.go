package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"project/internal/model"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// mockService is a test double for ItemsService
type mockService struct {
	expectedIDs []string
	returnItems []model.ItemDTO
	calledWith  []string
}

func (m *mockService) GetItemsByIDs(ctx context.Context, ids []string) []model.ItemDTO {
	m.calledWith = ids
	return m.returnItems
}

func TestGetData(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		body           interface{}
		mockReturn     []model.ItemDTO
		expectedStatus int
		expectedBody   string // substring check
	}{
		{
			name:           "invalid JSON body",
			body:           "{invalid_json}",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body",
		},
		{
			name: "empty IDs returns bad request",
			body: RequestBody{
				ItemID: []ItemData{},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid request body",
		},
		{
			name: "no items found returns 404",
			body: RequestBody{
				ItemID: []ItemData{{ID: "123"}},
			},
			mockReturn:     []model.ItemDTO{},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "no items found",
		},
		{
			name: "successful request returns 200 with data payload",
			body: RequestBody{
				ItemID: []ItemData{{ID: "1"}},
			},
			mockReturn: []model.ItemDTO{
				{ID: "1", Name: "Item One"},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "\"data\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := &mockService{returnItems: tt.mockReturn}
			h := NewHandler(mockSvc)

			// Create test server
			router := gin.New()
			router.POST("/v1/api/items", h.GetData)

			var reqBody []byte
			var err error

			switch v := tt.body.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, err = json.Marshal(v)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
			}

			req, _ := http.NewRequest(http.MethodPost, "/v1/api/items", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			if !bytes.Contains(rec.Body.Bytes(), []byte(tt.expectedBody)) {
				t.Errorf("expected response body to contain %q, got %s", tt.expectedBody, rec.Body.String())
			}

			// For successful request, ensure correct IDs were passed to service
			if tt.expectedStatus == http.StatusOK {
				expectedIDs := []string{"1"}
				assert.Equal(t, expectedIDs, mockSvc.calledWith, "service should be called with correct IDs")
			}
		})
	}
}
