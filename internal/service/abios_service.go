package service

import (
	"context"

	"github.com/benjaminmishra/abios-apis/internal/abios"
	models "github.com/benjaminmishra/abios-apis/internal/models"
)

type LiveService interface {
	GetLiveSeries(ctx context.Context) ([]models.SeriesDetails, error)
	GetLivePlayers(ctx context.Context) ([]models.Player, error)
	GetLiveTeams(ctx context.Context) ([]models.Team, error)
}

type abiosLiveService struct {
	client abios.AbiosClient
}

func NewAbiosLiveService(client abios.AbiosClient) *abiosLiveService {
	return &abiosLiveService{client: client}
}

func (s *abiosLiveService) GetLiveSeries(ctx context.Context) ([]models.SeriesDetails, error) {
	series, err := s.client.GetLiveSeries(ctx)
	if err != nil {
		return nil, err
	}

	if len(series) == 0 {
		return nil, nil
	}

	result := make([]models.SeriesDetails, len(series))
	for i, sr := range series {
		result[i] = models.SeriesDetails{
			ID:    sr.ID,
			Title: sr.Title,
		}
	}

	return result, nil
}

func (s *abiosLiveService) GetLivePlayers(ctx context.Context) ([]models.Player, error) {
	series, err := s.client.GetLiveSeries(ctx)
	if err != nil {
		return nil, err
	}

	// collect roster IDs from live series
	liveRosterIDs := []int{}
	for _, sr := range series {
		for _, p := range sr.Participants {
			liveRosterIDs = append(liveRosterIDs, p.Roster.ID)
		}
	}

	liveRosters, err := s.client.GetRostersByID(ctx, liveRosterIDs)
	if err != nil {
		return nil, err
	}

	// collect unique player IDs from rosters (assuming we have duplicates)
	playerIDMap := map[int]struct{}{}
	for _, r := range liveRosters {
		for _, p := range r.LineUp.Players {
			playerIDMap[p.ID] = struct{}{}
		}
	}
	uniquePlayerIDs := mapKeysToSlice(playerIDMap)

	result, err := s.client.GetPlayersByID(ctx, uniquePlayerIDs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *abiosLiveService) GetLiveTeams(ctx context.Context) ([]models.Team, error) {
	series, err := s.client.GetLiveSeries(ctx)
	if err != nil {
		return nil, err
	}

	// collect roster IDs
	liveRosterIDs := []int{}
	for _, sr := range series {
		for _, p := range sr.Participants {
			liveRosterIDs = append(liveRosterIDs, p.Roster.ID)
		}
	}

	liveRosters, err := s.client.GetRostersByID(ctx, liveRosterIDs)
	if err != nil {
		return nil, err
	}

	teamIDMap := map[int]struct{}{}
	for _, r := range liveRosters {
		teamIDMap[r.TeamId.ID] = struct{}{}
	}

	uniqueTeamIDs := mapKeysToSlice(teamIDMap)

	result, err := s.client.GetTeamsByID(ctx, uniqueTeamIDs)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func mapKeysToSlice(m map[int]struct{}) []int {
	out := make([]int, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
