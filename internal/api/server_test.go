package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/benjaminmishra/abios-apis/internal/api"
	"github.com/benjaminmishra/abios-apis/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLiveService struct {
	mock.Mock
}

func (m *mockLiveService) GetLiveSeries(ctx context.Context) ([]models.SeriesDetails, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.SeriesDetails), args.Error(1)
}

func (m *mockLiveService) GetLivePlayers(ctx context.Context) ([]models.Player, error) {
	args := m.Called(ctx)

	if s, ok := args.Get(0).([]models.Player); ok {
		return s, args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *mockLiveService) GetLiveTeams(ctx context.Context) ([]models.Team, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Team), args.Error(1)
}

func TestGetLiveSeries(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *mockLiveService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			setupMock: func(m *mockLiveService) {
				m.On("GetLiveSeries", mock.Anything).Return([]models.SeriesDetails{
					{ID: 1, Title: "Series 1"},
					{ID: 2, Title: "Series 2"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"title":"Series 1"},{"id":2,"title":"Series 2"}]`,
		},
		{
			name: "No Data",
			setupMock: func(m *mockLiveService) {
				m.On("GetLiveSeries", mock.Anything).Return([]models.SeriesDetails{}, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "No live series found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockService := new(mockLiveService)
			h := api.NewHandler(context.Background(), mockService)

			tt.setupMock(mockService)

			req := httptest.NewRequest(http.MethodGet, "/live/series", nil)
			w := httptest.NewRecorder()

			h.GetLiveSeries(w, req)

			result := w.Result()
			assert.Equal(t, tt.expectedStatus, result.StatusCode)

			body := w.Body.String()
			if result.StatusCode == http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, body)
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetLivePlayers(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *mockLiveService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			setupMock: func(m *mockLiveService) {
				m.On("GetLivePlayers", mock.Anything).Return([]models.Player{
					{ID: 1, Nickname: "Player 1"},
					{ID: 2, Nickname: "Player 2"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"nick_name":"Player 1"},{"id":2,"nick_name":"Player 2"}]`,
		},
		{
			name: "No Data",
			setupMock: func(m *mockLiveService) {
				m.On("GetLivePlayers", mock.Anything).Return([]models.Player{}, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "No live players found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockService := new(mockLiveService)
			h := api.NewHandler(context.Background(), mockService)

			tt.setupMock(mockService)

			req := httptest.NewRequest(http.MethodGet, "/live/players", nil)
			w := httptest.NewRecorder()

			h.GetLivePlayers(w, req)

			result := w.Result()
			assert.Equal(t, tt.expectedStatus, result.StatusCode)

			body := w.Body.String()
			if result.StatusCode == http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, body)
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestGetLiveTeams(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(m *mockLiveService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			setupMock: func(m *mockLiveService) {
				m.On("GetLiveTeams", mock.Anything).Return([]models.Team{
					{ID: 1, Name: "Team 1"},
					{ID: 2, Name: "Team 2"},
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `[{"id":1,"name":"Team 1"},{"id":2,"name":"Team 2"}]`,
		},
		{
			name: "No Data",
			setupMock: func(m *mockLiveService) {
				m.On("GetLiveTeams", mock.Anything).Return([]models.Team{}, nil)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   "No live teams found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockService := new(mockLiveService)
			h := api.NewHandler(context.Background(), mockService)

			tt.setupMock(mockService)

			req := httptest.NewRequest(http.MethodGet, "/live/teams", nil)
			w := httptest.NewRecorder()

			h.GetLiveTeams(w, req)

			result := w.Result()
			assert.Equal(t, tt.expectedStatus, result.StatusCode)

			body := w.Body.String()
			if result.StatusCode == http.StatusOK {
				assert.JSONEq(t, tt.expectedBody, body)
			} else {
				assert.Equal(t, tt.expectedBody, body)
			}

			mockService.AssertExpectations(t)
		})
	}
}
