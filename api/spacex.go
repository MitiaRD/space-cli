package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/MitiaRD/ReMarkable-cli/model"
)

func GetAllCrewMembers() (map[string]model.Crew, error) {
	url := "https://api.spacexdata.com/v4/crew"
	crew, err := fetchFromAPI[[]model.Crew](url)
	if err != nil {
		return nil, err
	}
	crewMap := make(map[string]model.Crew)
	for _, crewMember := range crew {
		crewMap[crewMember.ID] = crewMember
	}
	return crewMap, nil
}

func GetRocket(rocketId string) (model.Rocket, error) {
	url := fmt.Sprintf("https://api.spacexdata.com/v4/rockets/%s", rocketId)
	return fetchFromAPI[model.Rocket](url)
}

func GetLaunchpad(launchpadId string) (model.Launchpad, error) {
	url := fmt.Sprintf("https://api.spacexdata.com/v4/launchpads/%s", launchpadId)
	return fetchFromAPI[model.Launchpad](url)
}

func fetchFromAPI[T any](url string) (T, error) {
	var result T
	maxRetries := 5
	baseDelay := 200 * time.Millisecond
	maxDelay := 5 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := http.Get(url)
		if err != nil {
			if attempt == maxRetries {
				return result, fmt.Errorf("HTTP request failed after %d attempts: %v", maxRetries+1, err)
			}
			delay := calculateBackoffDelay(attempt, baseDelay, maxDelay)
			fmt.Printf("Request failed, retrying in %v (attempt %d/%d): %v\n", delay, attempt+1, maxRetries+1, err)
			time.Sleep(delay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 429 || (resp.StatusCode >= 500 && resp.StatusCode < 600) {
			if attempt == maxRetries {
				body, _ := io.ReadAll(resp.Body)
				return result, fmt.Errorf("API returned status %d after %d attempts: %s", resp.StatusCode, maxRetries+1, string(body))
			}

			retryAfter := resp.Header.Get("Retry-After")
			var delay time.Duration
			if retryAfter != "" {
				if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
					delay = seconds
				} else {
					delay = calculateBackoffDelay(attempt, baseDelay, maxDelay)
				}
			} else {
				delay = calculateBackoffDelay(attempt, baseDelay, maxDelay)
			}

			fmt.Printf("Rate limited (429) or server error (%d), retrying in %v (attempt %d/%d)\n", resp.StatusCode, delay, attempt+1, maxRetries+1)
			time.Sleep(delay)
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return result, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			if attempt == maxRetries {
				return result, fmt.Errorf("failed to read response body after %d attempts: %v", maxRetries+1, err)
			}
			delay := calculateBackoffDelay(attempt, baseDelay, maxDelay)
			fmt.Printf("Failed to read response, retrying in %v (attempt %d/%d): %v\n", delay, attempt+1, maxRetries+1, err)
			time.Sleep(delay)
			continue
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			if attempt == maxRetries {
				return result, fmt.Errorf("failed to parse JSON after %d attempts: %v. Response: %s", maxRetries+1, err, string(body[:min(200, len(body))]))
			}
			delay := calculateBackoffDelay(attempt, baseDelay, maxDelay)
			fmt.Printf("Failed to parse JSON, retrying in %v (attempt %d/%d): %v\n", delay, attempt+1, maxRetries+1, err)
			time.Sleep(delay)
			continue
		}

		return result, nil
	}

	return result, fmt.Errorf("unexpected error: max retries exceeded")
}

func calculateBackoffDelay(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	delay := float64(baseDelay) * float64(attempt+1)

	jitter := delay * 0.25 * (rand.Float64()*2 - 1)
	delay += jitter

	if delay > float64(maxDelay) {
		delay = float64(maxDelay)
	}

	return time.Duration(delay)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetLaunchesWithQuery(query map[string]interface{}) ([]model.Launch, error) {
	return fetchLaunchesWithQuery("https://api.spacexdata.com/v4/launches/query", query)
}

func fetchLaunchesWithQuery(url string, query map[string]interface{}) ([]model.Launch, error) {
	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	maxRetries := 5
	baseDelay := 200 * time.Millisecond
	maxDelay := 5 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, fmt.Errorf("HTTP request failed after %d attempts: %v", maxRetries+1, err)
			}
			delay := calculateBackoffDelay(attempt, baseDelay, maxDelay)
			fmt.Printf("Request failed, retrying in %v (attempt %d/%d): %v\n", delay, attempt+1, maxRetries+1, err)
			time.Sleep(delay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == 429 || (resp.StatusCode >= 500 && resp.StatusCode < 600) {
			if attempt == maxRetries {
				body, _ := io.ReadAll(resp.Body)
				return nil, fmt.Errorf("API returned status %d after %d attempts: %s", resp.StatusCode, maxRetries+1, string(body))
			}

			retryAfter := resp.Header.Get("Retry-After")
			var delay time.Duration
			if retryAfter != "" {
				if seconds, err := time.ParseDuration(retryAfter + "s"); err == nil {
					delay = seconds
				} else {
					delay = calculateBackoffDelay(attempt, baseDelay, maxDelay)
				}
			} else {
				delay = calculateBackoffDelay(attempt, baseDelay, maxDelay)
			}

			fmt.Printf("Rate limited (429) or server error (%d), retrying in %v (attempt %d/%d)\n", resp.StatusCode, delay, attempt+1, maxRetries+1)
			time.Sleep(delay)
			continue
		}

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			if attempt == maxRetries {
				return nil, fmt.Errorf("failed to read response body after %d attempts: %v", maxRetries+1, err)
			}
			delay := calculateBackoffDelay(attempt, baseDelay, maxDelay)
			fmt.Printf("Failed to read response, retrying in %v (attempt %d/%d): %v\n", delay, attempt+1, maxRetries+1, err)
			time.Sleep(delay)
			continue
		}

		var result struct {
			Docs []model.Launch `json:"docs"`
		}

		err = json.Unmarshal(body, &result)
		if err != nil {
			if attempt == maxRetries {
				return nil, fmt.Errorf("failed to parse JSON after %d attempts: %v. Response: %s", maxRetries+1, err, string(body[:min(200, len(body))]))
			}
			delay := calculateBackoffDelay(attempt, baseDelay, maxDelay)
			fmt.Printf("Failed to parse JSON, retrying in %v (attempt %d/%d): %v\n", delay, attempt+1, maxRetries+1, err)
			time.Sleep(delay)
			continue
		}

		return result.Docs, nil
	}

	return nil, fmt.Errorf("unexpected error: max retries exceeded")
}
