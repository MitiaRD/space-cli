package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/MitiaRD/ReMarkable-cli/model"
)

func GetCrewMember(crewId string) (model.Crew, error) {
	url := fmt.Sprintf("https://api.spacexdata.com/v4/crew/%s", crewId)
	return fetchFromAPI[model.Crew](url)
}

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

// SpaceX api does not seem to support limit or order by date flags...
func GetPastLaunches(limit int) ([]model.Launch, error) {
	url := "https://api.spacexdata.com/v4/launches/past"
	launches, err := fetchFromAPI[[]model.Launch](url)
	if err != nil {
		return nil, err
	}

	sort.Slice(launches, func(i, j int) bool {
		return launches[i].Date.After(launches[j].Date)
	})

	if limit > 0 && limit < len(launches) {
		return launches[:limit], nil
	}
	return launches, nil
}

func GetUpcomingLaunches(limit int) ([]model.Launch, error) {
	url := "https://api.spacexdata.com/v4/launches/upcoming"
	launches, err := fetchFromAPI[[]model.Launch](url)
	if err != nil {
		return nil, err
	}

	sort.Slice(launches, func(i, j int) bool {
		return launches[i].Date.Before(launches[j].Date)
	})

	if limit > 0 && limit < len(launches) {
		return launches[:limit], nil
	}
	return launches, nil
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
