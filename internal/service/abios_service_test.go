package service_test

import (
	"context"
	"testing"

	"github.com/benjaminmishra/abios-apis/internal/models"
	"github.com/benjaminmishra/abios-apis/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAbiosClient struct {
	mock.Mock
}

func (m *mockAbiosClient) GetLiveSeries(ctx context.Context) ([]models.Series, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Series), args.Error(1)
}

func (m *mockAbiosClient) GetRostersByID(ctx context.Context, ids []int) ([]models.Roster, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]models.Roster), args.Error(1)
}

func (m *mockAbiosClient) GetPlayersByID(ctx context.Context, ids []int) ([]models.Player, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]models.Player), args.Error(1)
}

func (m *mockAbiosClient) GetTeamsByID(ctx context.Context, ids []int) ([]models.Team, error) {
	args := m.Called(ctx, ids)
	return args.Get(0).([]models.Team), args.Error(1)
}

func TestNewAbiosLiveService(t *testing.T) {
	mockClient := new(mockAbiosClient)
	service := service.NewAbiosLiveService(mockClient)
	assert.NotNil(t, service)
}

func TestGetLiveSeries(t *testing.T) {
	mockClient := new(mockAbiosClient)
	service := service.NewAbiosLiveService(mockClient)
	ctx := context.Background()

	mockSeries := []models.Series{
		{
			ID:    1,
			Title: "Series 1",
		},
		{
			ID:    2,
			Title: "Series 2",
		},
	}

	mockClient.On("GetLiveSeries", ctx).Return(mockSeries, nil)

	result, err := service.GetLiveSeries(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, mockSeries[0].ID, result[0].ID)
	assert.Equal(t, mockSeries[0].Title, result[0].Title)

	mockClient.AssertExpectations(t)
}

func TestGetLivePlayers(t *testing.T) {
	mockClient := new(mockAbiosClient)
	service := service.NewAbiosLiveService(mockClient)
	ctx := context.Background()

	mockSeries := []models.Series{
		{
			ID:    1,
			Title: "Series 1",
			Participants: []models.Participant{
				{
					Roster: models.Roster{
						ID: 10,
					},
				},
				{
					Roster: models.Roster{
						ID: 20,
					},
				},
			},
		},
	}

	mockRosters := []models.Roster{
		{
			ID: 10,
			TeamId: models.TeamId{
				ID: 1,
			},
			LineUp: models.LineUp{
				Players: []models.PlayerId{
					{ID: 100},
					{ID: 101},
				},
			},
		},
	}

	mockClient.On("GetLiveSeries", ctx).Return(mockSeries, nil)
	mockClient.On("GetRostersByID", ctx, []int{10, 20}).Return(mockRosters, nil)
	mockClient.On("GetPlayersByID", ctx, []int{100, 101}).Return([]models.Player{
		{
			ID:       100,
			Nickname: "Player A",
		},
		{
			ID:       101,
			Nickname: "Player B",
		},
	}, nil)

	result, err := service.GetLivePlayers(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockClient.AssertExpectations(t)
}

func TestGetLiveTeams(t *testing.T) {
	mockClient := new(mockAbiosClient)
	s := service.NewAbiosLiveService(mockClient) // Renamed service to s
	ctx := context.Background()

	mockSeries := []models.Series{
		{
			ID:    1,
			Title: "Series 1",
			Participants: []models.Participant{
				{Roster: models.Roster{ID: 10}},
				{Roster: models.Roster{ID: 20}},
			},
		},
	}

	mockRosters := []models.Roster{
		{
			ID:     10,
			TeamId: models.TeamId{ID: 100},
		},
		{
			ID:     20,
			TeamId: models.TeamId{ID: 200},
		},
	}

	mockClient.On("GetLiveSeries", ctx).Return(mockSeries, nil)
	mockClient.On("GetRostersByID", ctx, []int{10, 20}).Return(mockRosters, nil)

	expectedTeamIDs := []int{100, 200}
	mockClient.On("GetTeamsByID", ctx, mock.MatchedBy(func(ids []int) bool {
		return assert.ElementsMatch(t, expectedTeamIDs, ids)
	})).Return([]models.Team{
		{ID: 100, Name: "Team A"},
		{ID: 200, Name: "Team B"},
	}, nil)

	result, err := s.GetLiveTeams(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	mockClient.AssertExpectations(t)
}
