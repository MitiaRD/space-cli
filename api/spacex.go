package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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

func GetLaunchesWithQuery(query map[string]interface{}) ([]model.Launch, error) {
	return fetchLaunchesWithQuery("https://api.spacexdata.com/v4/launches/query", query)
}

func fetchLaunchesWithQuery(url string, query map[string]interface{}) ([]model.Launch, error) {
	jsonData, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Docs []model.Launch `json:"docs"`
	}

	err = json.Unmarshal(body, &result)

	return result.Docs, err
}
