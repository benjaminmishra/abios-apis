package abios

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/benjaminmishra/abios-apis/internal/models"
	"golang.org/x/time/rate"
)

type AbiosClient interface {
	GetLiveSeries(ctx context.Context) ([]models.Series, error)
	GetRostersByID(ctx context.Context, rosterIDs []int) ([]models.Roster, error)
	GetTeamsByID(ctx context.Context, teamIDs []int) ([]models.Team, error)
	GetPlayersByID(ctx context.Context, playerIDs []int) ([]models.Player, error)
}

type client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

func NewClient(baseURL, token string, requestTimeoutSec int, requestsPerSec, burst int) AbiosClient {

	transport := &authTransport{
		token: token,
		transport: &rateLimitTransport{
			limiter: rate.NewLimiter(rate.Limit(requestsPerSec), burst),
			transport: &retryTransport{
				transport:  http.DefaultTransport,
				maxRetries: 3,
			},
		},
	}

	return &client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout:   time.Duration(requestTimeoutSec) * time.Second,
			Transport: transport,
		},
	}
}

func (c *client) GetLiveSeries(ctx context.Context) ([]models.Series, error) {
	endpoint := fmt.Sprintf("%s/series?filter=lifecycle=live", c.baseURL)

	return getAndDecode[models.Series](ctx, c.httpClient, endpoint)
}

func (c *client) GetRostersByID(ctx context.Context, rosterIDs []int) ([]models.Roster, error) {
	filter := buildIDFilter("id", rosterIDs)

	params := url.Values{}
	params.Add("filter", filter)

	endpoint := fmt.Sprintf("%s/rosters?%s", c.baseURL, params.Encode())

	return getAndDecode[models.Roster](ctx, c.httpClient, endpoint)
}

func (c *client) GetTeamsByID(ctx context.Context, teamIDs []int) ([]models.Team, error) {
	filter := buildIDFilter("id", teamIDs)

	params := url.Values{}
	params.Add("filter", filter)

	endpoint := fmt.Sprintf("%s/teams?%s", c.baseURL, params.Encode())

	return getAndDecode[models.Team](ctx, c.httpClient, endpoint)
}

func (c *client) GetPlayersByID(ctx context.Context, playerIDs []int) ([]models.Player, error) {
	filter := buildIDFilter("id", playerIDs)

	params := url.Values{}
	params.Add("filter", filter)

	endpoint := fmt.Sprintf("%s/players?%s", c.baseURL, params.Encode())

	return getAndDecode[models.Player](ctx, c.httpClient, endpoint)
}

// Helpers
func buildIDFilter(key string, ids []int) string {
	strIDs := make([]string, len(ids))
	for i, id := range ids {
		strIDs[i] = fmt.Sprintf("%d", id)
	}
	return fmt.Sprintf("%s<={%s}", key, strings.Join(strIDs, ","))
}

func parseRetryAfter(header string) time.Duration {
	if header == "" {
		return time.Second
	}
	if secs, err := strconv.Atoi(header); err == nil {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(header); err == nil {
		return time.Until(t)
	}
	return time.Second
}

func getAndDecode[T any](ctx context.Context, httpClient *http.Client, endpoint string) ([]T, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("abios: unexpected status %d", resp.StatusCode)
	}

	var result []T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
