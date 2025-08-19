package api

import (
	"fmt"
	"os"
	"time"

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

func GetAsteriods(queryParams string) (model.NasaAsteriod, error) {
	url := "https://api.nasa.gov/neo/rest/v1/feed" + queryParams + "&api_key=" + GetNASAAPIKey()

	asteriods, err := fetchFromAPI[model.NasaAsteriod](url)
	if err != nil {
		return model.NasaAsteriod{}, err
	}

	return asteriods, nil
}

func BuildWeatherEventsQueryParams(long, lat float64, date time.Time) string {
	return fmt.Sprintf("?bbox=%f,%f,%f,%f&start=%s&end=%s", long-1, lat-1, long+1, lat+1, date.Format("2006-01-02"), date.Format("2006-01-02"))
}
