/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/MitiaRD/ReMarkable-cli/api"
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
		query := buildLaunchQuery(cmd)

		launches, err := api.GetLaunchesWithQuery(query)
		if err != nil {
			fmt.Printf("Error fetching launches: %v\n", err)
			return
		}

		fmt.Printf("\n🚀 Launches (showing %d):\n", len(launches))
		fmt.Println(strings.Repeat("-", 80))

		for _, launch := range launches {
			rocket, err := api.GetRocket(launch.RocketId)
			if err != nil {
				fmt.Printf("Error fetching rocket: %v\n", err)
				continue
			}

			status := "❌ Failed"
			if launch.Success {
				status = "✅ Success"
			}

			fmt.Printf("📅 %s\n", launch.Date.Format("2006-01-02 15:04"))
			fmt.Printf("   🏷️  %s\n", launch.Name)
			fmt.Printf("   %s\n", status)
			fmt.Printf("   🚀 %s\n", rocket.Name)
			fmt.Printf("   💰 $%v million\n", rocket.CostPerLaunch)

			if len(launch.Crew) > 0 {
				crewMap, err := api.GetAllCrewMembers()
				if err == nil {
					fmt.Printf("   👥 Crew: ")
					crewNames := []string{}
					for _, crewId := range launch.Crew {
						if crew, exists := crewMap[crewId]; exists {
							crewNames = append(crewNames, crew.Name)
						}
					}
					fmt.Printf("%s\n", strings.Join(crewNames, ", "))
				}
			}

			fmt.Println()
		}
	},
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
	}

	limit, _ := cmd.Flags().GetInt("limit")
	if limit > 0 {
		query["options"].(map[string]interface{})["limit"] = limit
	}

	return query
}

func init() {
	rootCmd.AddCommand(launchesCmd)

	launchesCmd.Flags().IntP("limit", "l", 10, "Number of past launches to show")
	launchesCmd.Flags().StringP("start", "s", "", "Start date (YYYY-MM-DD)")
	launchesCmd.Flags().StringP("end", "e", "", "End date (YYYY-MM-DD)")
	launchesCmd.Flags().BoolP("failed", "f", true, "Filter for successful launches only")
	launchesCmd.Flags().BoolP("upcoming", "c", false, "Filter for upcoming launches only")
}
