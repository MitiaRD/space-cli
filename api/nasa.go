package api

import (
	"os"

	"github.com/MitiaRD/ReMarkable-cli/model"
)

func GetNASAAPIKey() string {
	return os.Getenv("NASA_API_KEY")
}

func GetEarthEvents(queryParams string) ([]model.NasaEarthEvent, error) {
	url := "https://eonet.gsfc.nasa.gov/api/v3/events" + queryParams

	events, err := fetchFromAPI[model.NasaEarth](url)
	if err != nil {
		return []model.NasaEarthEvent{}, err
	}

	return events.Events, nil
}
