/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/MitiaRD/ReMarkable-cli/model"
	"github.com/spf13/cobra"
)

var launchesCmd = &cobra.Command{
	Use:   "launches",
	Short: "Explore space launch information",
	Long: `Launches provides comprehensive information about space launches.
	
Available subcommands:
  limit        - Limit the number of launches to show,
  start        - Start date (YYYY-MM-DD),
  end          - End date (YYYY-MM-DD),
  failed       - Filter for failed launches only,
  upcoming     - Filter for upcoming launches only`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		config, err := LoadConfiguration()
		if err != nil {
			fmt.Printf("Error loading configuration: %v\n", err)
			return
		}

		logger := SetupLogger()
		service := NewLaunchesService(config, logger)

		query := buildLaunchQuery(cmd)

		launches, err := service.GetLaunches(ctx, query)
		if err != nil {
			logger.Error("failed to fetch launches", "error", err)
			fmt.Printf("Error fetching launches: %v\n", err)
			return
		}

		fmt.Printf("\n🚀 Launches (showing %d):\n", len(launches))

		rockets, err := service.GetRockets(ctx)
		if err != nil {
			logger.Error("failed to fetch rockets", "error", err)
		}

		cost, _ := cmd.Flags().GetBool("cost")
		if cost {
			getCosts(launches, rockets)
			return
		}

		crewMap, err := service.GetCrewMembers(ctx)
		if err != nil {
			logger.Error("failed to fetch crew members", "error", err)
		}

		launchpads, err := service.GetLaunchpads(ctx)
		if err != nil {
			logger.Error("failed to fetch launchpads", "error", err)
		}

		fmt.Println(strings.Repeat("-", 80))
		for _, launch := range launches {

			status := "❓ Unknown"
			if launch.Success != nil {
				if *launch.Success {
					status = "✅ Success"
				} else {
					status = "❌ Failed"
				}
			}

			fmt.Printf("📅 %s\n", launch.Date.Format("2006-01-02 15:04"))
			fmt.Printf("   🏷️  %s\n", launch.Name)
			fmt.Printf("   %s\n", status)
			fmt.Printf("   🚀 %s\n", rockets[launch.RocketId].Name)
			if launch.Details != "" {
				fmt.Printf("   ℹ️ %v \n", launch.Details)
			}

			if len(launch.Crew) > 0 && crewMap != nil {
				fmt.Printf("   👥 Crew: ")
				crewNames := []string{}
				for _, crewId := range launch.Crew {
					if crew, exists := crewMap[crewId]; exists {
						crewNames = append(crewNames, crew.Name)
					}
				}
				fmt.Printf("%s\n", strings.Join(crewNames, ", "))
			}

			launchpad, _ := cmd.Flags().GetBool("launchpad")
			if launchpad && launchpads != nil {
				launchpad, exists := launchpads[launch.LaunchpadId]
				if !exists {
					fmt.Printf("Launchpad not found for launch %s\n", launch.LaunchpadId)
					continue
				}

				fmt.Printf("   📍 %s\n", launchpad.Name)
				fmt.Printf("      (%s)\n", launchpad.Details)

				if weatherEvents, _ := cmd.Flags().GetBool("weather"); weatherEvents {
					weatherEvents, err := service.GetEarthEvents(ctx, launchpad.Longitude, launchpad.Latitude, launch.Date)
					if err != nil {
						logger.Error("failed to fetch weather events", "error", err)
						continue
					}
					if len(weatherEvents) == 0 {
						fmt.Printf("   🌤️  No warning events found from Nasa for this time & location\n")

					} else {
						for _, event := range weatherEvents {
							fmt.Printf("    🌤️  %s (%s)\n", event.Title, event.Description)
						}
					}
				}
			}

			if asteroids, _ := cmd.Flags().GetBool("asteroids"); asteroids {
				asteroids, err := service.GetAsteroids(ctx, launch.Date)
				if err != nil {
					logger.Error("failed to fetch asteroids", "error", err)
				}
				hazardous := 0
				nonHazardous := 0
				maxDiameter := 0.0
				minDiameter := 0.0
				for _, asteroids := range asteroids.NearEarthObjects {
					for _, asteroid := range asteroids {
						if asteroid.Hazardous {
							hazardous++
						} else {
							nonHazardous++
						}
						if asteroid.Diameter.Meters.Estimated > maxDiameter {
							maxDiameter = asteroid.Diameter.Meters.Estimated
						}
						if asteroid.Diameter.Meters.Estimated < minDiameter || minDiameter == 0 {
							minDiameter = asteroid.Diameter.Meters.Estimated
						}
					}
				}
				fmt.Printf("   🌍  total number of near earth asteroids %d (hazardous: %d, non-hazardous: %d) with diameters ranging from %f to %f meters\n", asteroids.ElementCount, hazardous, nonHazardous, minDiameter, maxDiameter)
			}

			fmt.Println()
		}
	},
}

func getCosts(launches []model.Launch, rockets map[string]model.Rocket) (int, error) {
	totalCost := 0
	var wg sync.WaitGroup

	costChan := make(chan int, len(launches))

	for _, launch := range launches {
		wg.Add(1)
		go func(launch model.Launch) {
			defer wg.Done()
			if rocket, exists := rockets[launch.RocketId]; exists {
				costChan <- rocket.CostPerLaunch
			} else {
				costChan <- 0
			}
		}(launch)
	}

	go func() {
		wg.Wait()
		close(costChan)
	}()

	for cost := range costChan {
		totalCost += cost
	}

	fmt.Printf("Total cost: $%d million\n", totalCost)
	return totalCost, nil
}

func buildAsteroidsQueryParams(date time.Time) string {
	return fmt.Sprintf("?start_date=%s&end_date=%s", date.Format("2006-01-02"), date.Format("2006-01-02"))
}

func buildLaunchQuery(cmd *cobra.Command) map[string]interface{} {
	query := map[string]interface{}{
		"query": map[string]interface{}{},
		"options": map[string]interface{}{
			"sort": map[string]interface{}{
				"date_utc": "desc",
			},
		},
	}

	startDate, _ := cmd.Flags().GetString("start")
	endDate, _ := cmd.Flags().GetString("end")

	if startDate != "" && endDate != "" {
		dateQuery := map[string]interface{}{
			"$gte": startDate + "T00:00:00.000Z",
			"$lte": endDate + "T23:59:59.999Z",
		}
		query["query"].(map[string]interface{})["date_utc"] = dateQuery
	}

	failed, _ := cmd.Flags().GetBool("failed")
	if failed {
		query["query"].(map[string]interface{})["success"] = false
	}

	upcoming, _ := cmd.Flags().GetBool("upcoming")
	if upcoming {
		query["query"].(map[string]interface{})["upcoming"] = true
	} else {
		query["query"].(map[string]interface{})["upcoming"] = false
	}

	limit, _ := cmd.Flags().GetInt("limit")
	if limit > 0 {
		query["options"].(map[string]interface{})["limit"] = limit
	}
	return query
}

func init() {
	rootCmd.AddCommand(launchesCmd)

	launchesCmd.Flags().IntP("limit", "l", 200, "Number of past launches to show")
	launchesCmd.Flags().StringP("start", "s", "", "Start date (YYYY-MM-DD)")
	launchesCmd.Flags().StringP("end", "e", "", "End date (YYYY-MM-DD)")
	launchesCmd.Flags().BoolP("failed", "f", false, "Filter for failed launches only")
	launchesCmd.Flags().BoolP("upcoming", "u", false, "Filter for upcoming launches only")
	launchesCmd.Flags().BoolP("cost", "c", false, "Get the total cost for all matching launches")
	launchesCmd.Flags().BoolP("launchpad", "p", false, "Show launchpad information")
	launchesCmd.Flags().BoolP("weather", "w", false, "Show launchpad location weather warning information")
	launchesCmd.Flags().BoolP("asteroids", "a", false, "Show near Earth orbiting asteroid information")
}
