package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stipochka/web_service/internal/mocks"
	"github.com/stipochka/web_service/internal/models"
	"github.com/stipochka/web_service/internal/service"
	"github.com/stretchr/testify/assert"
)

type Response struct {
	Data []models.Record `json:"data"`
}

func TestGetAllRecords(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := slog.Default()
	tests := []struct {
		name               string
		mockSetup          func(*mocks.GrpcClient)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name: "success",
			mockSetup: func(m *mocks.GrpcClient) {
				expectedRows := []models.Record{
					{ID: 1, Data: "temperature: test, humidity: test"},
					{ID: 2, Data: "temperature: test, humidity: test"},
				}
				m.On("GetAllRecords", context.Background()).Return(expectedRows, nil)
			},
			expectedStatusCode: 200,
			expectedResponse: Response{
				Data: []models.Record{
					{ID: 1, Data: "temperature: test, humidity: test"},
					{ID: 2, Data: "temperature: test, humidity: test"},
				},
			},
		},
		{
			name: "internal error",
			mockSetup: func(m *mocks.GrpcClient) {
				m.On("GetAllRecords", context.Background()).Return([]models.Record{}, errors.New("Internal error"))
			},
			expectedStatusCode: 500,
			expectedResponse: map[string]interface{}{
				"message": "Internal error",
			},
		},
		{
			name: "empty response",
			mockSetup: func(m *mocks.GrpcClient) {
				m.On("GetAllRecords", context.Background()).Return([]models.Record{}, nil)
			},
			expectedStatusCode: 200,
			expectedResponse:   Response{Data: []models.Record{}},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockService := new(mocks.GrpcClient)
			tc.mockSetup(mockService)

			h := NewHandler(mockLogger, &service.Service{GrpcClient: mockService})

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			h.getAllRecords(c)

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			if tc.expectedStatusCode == http.StatusInternalServerError {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Internal error", response["message"])
			} else {
				var response Response
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, response)
			}

			mockService.AssertExpectations(t)

		})
	}
}

func TestGetRecordByID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockLogger := slog.Default()
	tests := []struct {
		name               string
		urlIDParam         string
		mockSetup          func(*mocks.GrpcClient)
		expectedStatusCode int
		expectedResponse   interface{}
	}{
		{
			name:       "success with id 1",
			urlIDParam: "1",
			mockSetup: func(m *mocks.GrpcClient) {
				expectedRecord := models.Record{ID: 1, Data: "temperature: test, humidity: test"}
				m.On("GetRecordById", context.Background(), 1).Return(expectedRecord, nil)
			},
			expectedStatusCode: 200,
			expectedResponse:   models.Record{ID: 1, Data: "temperature: test, humidity: test"},
		},
		{
			name:       "empty result",
			urlIDParam: "4",
			mockSetup: func(m *mocks.GrpcClient) {
				expectedRecord := models.Record{}
				m.On("GetRecordById", context.Background(), 4).Return(expectedRecord, nil)
			},
			expectedStatusCode: 200,
			expectedResponse:   models.Record{},
		},
		{
			name:               "invalid record id",
			urlIDParam:         "das",
			mockSetup:          func(m *mocks.GrpcClient) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   map[string]interface{}{"message": "invalid id"},
		},
		{
			name:       "internal error",
			urlIDParam: "1",
			mockSetup: func(m *mocks.GrpcClient) {
				m.On("GetRecordById", context.Background(), 1).Return(models.Record{}, errors.New("no such id"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   map[string]interface{}{"message": "internal error"},
		},
		{
			name:               "missing id",
			urlIDParam:         "",
			mockSetup:          func(m *mocks.GrpcClient) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   map[string]interface{}{"message": "not given id"},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockService := new(mocks.GrpcClient)
			tc.mockSetup(mockService)

			h := NewHandler(mockLogger, &service.Service{GrpcClient: mockService})

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			if tc.urlIDParam != "" {
				c.Params = append(c.Params, gin.Param{Key: "id", Value: tc.urlIDParam})
			}

			h.getRecordById(c)

			assert.Equal(t, tc.expectedStatusCode, w.Code)
			if tc.expectedStatusCode == http.StatusInternalServerError {
				var errResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errResponse)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, errResponse)
			} else if tc.expectedStatusCode == http.StatusBadRequest {
				var errResponse map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &errResponse)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, errResponse)
			} else {
				var testResponse models.Record
				err := json.Unmarshal(w.Body.Bytes(), &testResponse)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResponse, testResponse)
			}

			mockService.AssertExpectations(t)
		})
	}

}
