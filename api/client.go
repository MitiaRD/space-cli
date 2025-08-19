package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/MitiaRD/ReMarkable-cli/model"
)

type Client interface {
	GetLaunchesWithQuery(ctx context.Context, query map[string]interface{}) ([]model.Launch, error)
	GetAllRockets(ctx context.Context) (map[string]model.Rocket, error)
	GetAllCrewMembers(ctx context.Context) (map[string]model.Crew, error)
	GetAllLaunchpads(ctx context.Context) (map[string]model.Launchpad, error)
	GetEarthEvents(ctx context.Context, queryParams string) ([]model.NasaEarthEvent, error)
	GetAsteroids(ctx context.Context, queryParams string) (model.NasaAsteroid, error)
}

type SpaceXClient struct {
	httpClient *http.Client
	logger     *slog.Logger
	config     *model.Config
}

type NASAClient struct {
	httpClient *http.Client
	logger     *slog.Logger
	config     *model.Config
}

func NewSpaceXClient(config *model.Config, logger *slog.Logger) *SpaceXClient {
	return &SpaceXClient{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		logger: logger,
		config: config,
	}
}

func NewNASAClient(config *model.Config, logger *slog.Logger) *NASAClient {
	return &NASAClient{
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		logger: logger,
		config: config,
	}
}

func (c *SpaceXClient) GetLaunchesWithQuery(ctx context.Context, query map[string]interface{}) ([]model.Launch, error) {
	return c.fetchLaunchesWithQuery(ctx, "https://api.spacexdata.com/v4/launches/query", query)
}

func (c *SpaceXClient) GetAllRockets(ctx context.Context) (map[string]model.Rocket, error) {
	url := "https://api.spacexdata.com/v4/rockets"
	rockets, err := fetchFromAPI[[]model.Rocket](c, ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rockets: %w", err)
	}

	rocketMap := make(map[string]model.Rocket)
	for _, rocket := range rockets {
		rocketMap[rocket.ID] = rocket
	}
	return rocketMap, nil
}

func (c *SpaceXClient) GetAllCrewMembers(ctx context.Context) (map[string]model.Crew, error) {
	url := "https://api.spacexdata.com/v4/crew"
	crew, err := fetchFromAPI[[]model.Crew](c, ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch crew members: %w", err)
	}

	crewMap := make(map[string]model.Crew)
	for _, crewMember := range crew {
		crewMap[crewMember.ID] = crewMember
	}
	return crewMap, nil
}

func (c *SpaceXClient) GetAllLaunchpads(ctx context.Context) (map[string]model.Launchpad, error) {
	url := "https://api.spacexdata.com/v4/launchpads"
	launchpads, err := fetchFromAPI[[]model.Launchpad](c, ctx, url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch launchpads: %w", err)
	}

	launchpadMap := make(map[string]model.Launchpad)
	for _, launchpad := range launchpads {
		launchpadMap[launchpad.ID] = launchpad
	}
	return launchpadMap, nil
}

func (c *NASAClient) GetEarthEvents(ctx context.Context, queryParams string) ([]model.NasaEarthEvent, error) {
	url := "https://eonet.gsfc.nasa.gov/api/v3/events" + queryParams
	events, err := fetchFromAPINASA[model.NasaEarth](c, ctx, url)
	if err != nil {
		return []model.NasaEarthEvent{}, fmt.Errorf("failed to fetch Earth events: %w", err)
	}
	return events.Events, nil
}

func (c *NASAClient) GetAsteroids(ctx context.Context, queryParams string) (model.NasaAsteroid, error) {
	url := "https://api.nasa.gov/neo/rest/v1/feed" + queryParams + "&api_key=" + c.config.NASAAPIKey
	asteroids, err := fetchFromAPINASA[model.NasaAsteroid](c, ctx, url)
	if err != nil {
		return model.NasaAsteroid{}, fmt.Errorf("failed to fetch asteroid data: %w", err)
	}
	return asteroids, nil
}

func fetchFromAPI[T any](c *SpaceXClient, ctx context.Context, url string) (T, error) {
	var result T

	for attempt := 0; attempt <= c.config.Retries; attempt++ {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			if attempt == c.config.Retries {
				return result, fmt.Errorf("failed to create request after %d attempts: %w", c.config.Retries+1, err)
			}
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Warn("request creation failed, retrying", "attempt", attempt+1, "delay", delay, "error", err)
			time.Sleep(delay)
			continue
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if attempt == c.config.Retries {
				return result, fmt.Errorf("HTTP request failed after %d attempts: %w", c.config.Retries+1, err)
			}
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Warn("HTTP request failed, retrying", "attempt", attempt+1, "delay", delay, "error", err)
			time.Sleep(delay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 429 || (resp.StatusCode >= 500 && resp.StatusCode < 600) {
			if attempt == c.config.Retries {
				body, _ := io.ReadAll(resp.Body)
				return result, fmt.Errorf("API returned status %d after %d attempts: %s", resp.StatusCode, c.config.Retries+1, string(body))
			}

			retryAfter := resp.Header.Get("Retry-After")
			var delay time.Duration
			if retryAfter != "" {
				if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
					delay = seconds
				} else {
					delay = c.calculateBackoffDelay(attempt)
				}
			} else {
				delay = c.calculateBackoffDelay(attempt)
			}

			c.logger.Warn("rate limited or server error, retrying", "status", resp.StatusCode, "attempt", attempt+1, "delay", delay)
			time.Sleep(delay)
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return result, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			if attempt == c.config.Retries {
				return result, fmt.Errorf("failed to read response body after %d attempts: %w", c.config.Retries+1, err)
			}
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Warn("failed to read response, retrying", "attempt", attempt+1, "delay", delay, "error", err)
			time.Sleep(delay)
			continue
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			if attempt == c.config.Retries {
				return result, fmt.Errorf("failed to parse JSON after %d attempts: %w. Response: %s", c.config.Retries+1, err, string(body[:min(200, len(body))]))
			}
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Warn("failed to parse JSON, retrying", "attempt", attempt+1, "delay", delay, "error", err)
			time.Sleep(delay)
			continue
		}

		return result, nil
	}

	return result, fmt.Errorf("unexpected error: max retries exceeded")
}

func fetchFromAPINASA[T any](c *NASAClient, ctx context.Context, url string) (T, error) {
	var result T

	for attempt := 0; attempt <= c.config.Retries; attempt++ {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			if attempt == c.config.Retries {
				return result, fmt.Errorf("failed to create request after %d attempts: %w", c.config.Retries+1, err)
			}
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Warn("request creation failed, retrying", "attempt", attempt+1, "delay", delay, "error", err)
			time.Sleep(delay)
			continue
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if attempt == c.config.Retries {
				return result, fmt.Errorf("HTTP request failed after %d attempts: %w", c.config.Retries+1, err)
			}
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Warn("HTTP request failed, retrying", "attempt", attempt+1, "delay", delay, "error", err)
			time.Sleep(delay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 429 || (resp.StatusCode >= 500 && resp.StatusCode < 600) {
			if attempt == c.config.Retries {
				body, _ := io.ReadAll(resp.Body)
				return result, fmt.Errorf("API returned status %d after %d attempts: %s", resp.StatusCode, c.config.Retries+1, string(body))
			}

			retryAfter := resp.Header.Get("Retry-After")
			var delay time.Duration
			if retryAfter != "" {
				if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
					delay = seconds
				} else {
					delay = c.calculateBackoffDelay(attempt)
				}
			} else {
				delay = c.calculateBackoffDelay(attempt)
			}

			c.logger.Warn("rate limited or server error, retrying", "status", resp.StatusCode, "attempt", attempt+1, "delay", delay)
			time.Sleep(delay)
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return result, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			if attempt == c.config.Retries {
				return result, fmt.Errorf("failed to read response body after %d attempts: %w", c.config.Retries+1, err)
			}
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Warn("failed to read response, retrying", "attempt", attempt+1, "delay", delay, "error", err)
			time.Sleep(delay)
			continue
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			if attempt == c.config.Retries {
				return result, fmt.Errorf("failed to parse JSON after %d attempts: %w. Response: %s", c.config.Retries+1, err, string(body[:min(200, len(body))]))
			}
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Warn("failed to parse JSON, retrying", "attempt", attempt+1, "delay", delay, "error", err)
			time.Sleep(delay)
			continue
		}

		return result, nil
	}

	return result, fmt.Errorf("unexpected error: max retries exceeded")
}

func (c *SpaceXClient) fetchLaunchesWithQuery(ctx context.Context, url string, query map[string]interface{}) ([]model.Launch, error) {
	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	for attempt := 0; attempt <= c.config.Retries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(jsonData)))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			if attempt == c.config.Retries {
				return nil, fmt.Errorf("HTTP request failed after %d attempts: %w", c.config.Retries+1, err)
			}
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Warn("HTTP request failed, retrying", "attempt", attempt+1, "delay", delay, "error", err)
			time.Sleep(delay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 429 || (resp.StatusCode >= 500 && resp.StatusCode < 600) {
			if attempt == c.config.Retries {
				body, _ := io.ReadAll(resp.Body)
				return nil, fmt.Errorf("API returned status %d after %d attempts: %s", resp.StatusCode, c.config.Retries+1, string(body))
			}

			retryAfter := resp.Header.Get("Retry-After")
			var delay time.Duration
			if retryAfter != "" {
				if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
					delay = seconds
				} else {
					delay = c.calculateBackoffDelay(attempt)
				}
			} else {
				delay = c.calculateBackoffDelay(attempt)
			}

			c.logger.Warn("rate limited or server error, retrying", "status", resp.StatusCode, "attempt", attempt+1, "delay", delay)
			time.Sleep(delay)
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			if attempt == c.config.Retries {
				return nil, fmt.Errorf("failed to read response body after %d attempts: %w", c.config.Retries+1, err)
			}
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Warn("failed to read response, retrying", "attempt", attempt+1, "delay", delay, "error", err)
			time.Sleep(delay)
			continue
		}

		var result struct {
			Docs []model.Launch `json:"docs"`
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			if attempt == c.config.Retries {
				return nil, fmt.Errorf("failed to parse JSON after %d attempts: %w. Response: %s", c.config.Retries+1, err, string(body[:min(200, len(body))]))
			}
			delay := c.calculateBackoffDelay(attempt)
			c.logger.Warn("failed to parse JSON, retrying", "attempt", attempt+1, "delay", delay, "error", err)
			time.Sleep(delay)
			continue
		}

		return result.Docs, nil
	}

	return nil, fmt.Errorf("unexpected error: max retries exceeded")
}

func (c *SpaceXClient) calculateBackoffDelay(attempt int) time.Duration {
	delay := float64(c.config.BaseDelay) * float64(attempt+1)

	jitter := delay * 0.25 * (rand.Float64()*2 - 1)
	delay += jitter

	if delay > float64(c.config.MaxDelay) {
		delay = float64(c.config.MaxDelay)
	}

	return time.Duration(delay)
}

func (c *NASAClient) calculateBackoffDelay(attempt int) time.Duration {
	delay := float64(c.config.BaseDelay) * float64(attempt+1)

	jitter := delay * 0.25 * (rand.Float64()*2 - 1)
	delay += jitter

	if delay > float64(c.config.MaxDelay) {
		delay = float64(c.config.MaxDelay)
	}

	return time.Duration(delay)
}

func BuildWeatherEventsQueryParams(long, lat float64, date time.Time) string {
	return fmt.Sprintf("?bbox=%f,%f,%f,%f&start=%s&end=%s", long-1, lat-1, long+1, lat+1, date.Format("2006-01-02"), date.Format("2006-01-02"))
}
