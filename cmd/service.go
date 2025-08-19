package cmd

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/MitiaRD/ReMarkable-cli/api"
	"github.com/MitiaRD/ReMarkable-cli/model"
)

type LaunchesService struct {
	spaceXClient *api.SpaceXClient
	nasaClient   *api.NASAClient
	logger       *slog.Logger
	config       *model.Config
}

func NewLaunchesService(config *model.Config, logger *slog.Logger) *LaunchesService {
	return &LaunchesService{
		spaceXClient: api.NewSpaceXClient(config, logger),
		nasaClient:   api.NewNASAClient(config, logger),
		logger:       logger,
		config:       config,
	}
}

func (s *LaunchesService) GetLaunches(ctx context.Context, query map[string]interface{}) ([]model.Launch, error) {
	return s.spaceXClient.GetLaunchesWithQuery(ctx, query)
}

func (s *LaunchesService) GetRockets(ctx context.Context) (map[string]model.Rocket, error) {
	return s.spaceXClient.GetAllRockets(ctx)
}

func (s *LaunchesService) GetCrewMembers(ctx context.Context) (map[string]model.Crew, error) {
	return s.spaceXClient.GetAllCrewMembers(ctx)
}

func (s *LaunchesService) GetLaunchpads(ctx context.Context) (map[string]model.Launchpad, error) {
	return s.spaceXClient.GetAllLaunchpads(ctx)
}

func (s *LaunchesService) GetEarthEvents(ctx context.Context, longitude, latitude float64, date time.Time) ([]model.NasaEarthEvent, error) {
	queryParams := api.BuildWeatherEventsQueryParams(longitude, latitude, date)
	return s.nasaClient.GetEarthEvents(ctx, queryParams)
}

func (s *LaunchesService) GetAsteroids(ctx context.Context, date time.Time) (model.NasaAsteroid, error) {
	queryParams := buildAsteroidsQueryParams(date)
	return s.nasaClient.GetAsteroids(ctx, queryParams)
}

func LoadConfiguration() (*model.Config, error) {
	config := model.DefaultConfig()

	if nasaKey := os.Getenv("NASA_API_KEY"); nasaKey != "" {
		config.NASAAPIKey = nasaKey
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func SetupLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	return slog.New(handler)
}
