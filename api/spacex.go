package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/MitiaRD/ReMarkable-cli/model"
)

func GetUpcomingLaunches(limit int) ([]model.Launch, error) {
	url := fmt.Sprintf("https://api.spacexdata.com/v4/launches/upcoming?limit=%d", limit)
	return fetchFromAPI[[]model.Launch](url)
}

func GetRocket(rocketId string) (model.Rocket, error) {
	url := fmt.Sprintf("https://api.spacexdata.com/v4/rockets/%s", rocketId)
	return fetchFromAPI[model.Rocket](url)
}

func fetchFromAPI[T any](url string) (T, error) {
	var result T

	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(body, &result)
	return result, err
}
