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

// launchesCmd represents the launches command
var launchesCmd = &cobra.Command{
	Use:   "launches",
	Short: "Explore space launch information",
	Long: `Launches provides comprehensive information about space launches.
	
Available subcommands:
  upcoming     - Get upcoming launches`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Use 'launches --help' to see available subcommands")
	},
}

var upcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "Get upcoming launches",
	Long:  `Get a list of upcoming spacex launches`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		if limit <= 0 {
			limit = 10
		}

		launches, err := api.GetUpcomingLaunches(limit)
		if err != nil {
			fmt.Printf("Error fetching upcoming launches: %v\n", err)
			return
		}

		fmt.Printf("\n🚀 Upcoming Launches (showing %d):\n", len(launches))
		fmt.Println(strings.Repeat("-", 80))

		for _, launch := range launches {
			rocket, err := api.GetRocket(launch.RocketId)
			if err != nil {
				fmt.Printf("Error fetching rocket: %v\n", err)
				continue
			}
			fmt.Printf("📅 %s\n", launch.Date.Format("2006-01-02 15:04"))
			fmt.Printf("   🏷️  %s\n", launch.ID)
			fmt.Printf("   🎬  %s\n", launch.Name)
			fmt.Printf("   🚀 %s\n", rocket.Name)
			fmt.Printf("   💰 $%v\n", rocket.CostPerLaunch)
			fmt.Printf("   ✅ %t\n", launch.Success)
			fmt.Println()
		}
	},
}

var pastCmd = &cobra.Command{
	Use:   "past",
	Short: "Get past launches",
	Long:  `Get a list of past spacex launches`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetInt("limit")
		if limit <= 0 {
			limit = 10
		}

		launches, err := api.GetPastLaunches(limit)
		if err != nil {
			fmt.Printf("Error fetching past launches: %v\n", err)
			return
		}

		fmt.Printf("\n🚀 Past Launches (showing %d):\n", len(launches))
		fmt.Println(strings.Repeat("-", 80))
		crew, err := api.GetAllCrewMembers()
		if err != nil {
			fmt.Printf("Error fetching crew: %v\n", err)
		}

		for _, launch := range launches {
			rocket, err := api.GetRocket(launch.RocketId)
			if err != nil {
				fmt.Printf("Error fetching rocket: %v\n", err)
				continue
			}
			fmt.Printf("📅 %s\n", launch.Date.Format("2006-01-02 15:04"))
			fmt.Printf("   🎬  %s\n", launch.Name)
			fmt.Printf("   🚀 %s\n", rocket.Name)
			fmt.Printf("   💰 $%v\n", rocket.CostPerLaunch)
			fmt.Printf("   ✅ %t\n", launch.Success)
			fmt.Printf("   📂 %s\n", launch.Details)
			for _, crewId := range launch.Crew {
				crew, ok := crew[crewId]
				if !ok {
					fmt.Printf("   👤 %s\n", crewId)
				} else {
					fmt.Printf("   👤 %s\n", crew.Name)
				}
			}
			fmt.Println()
		}
	},
}

func init() {
	rootCmd.AddCommand(launchesCmd)

	// Add subcommands
	launchesCmd.AddCommand(upcomingCmd)
	launchesCmd.AddCommand(pastCmd)

	// Add flags for subcommands
	upcomingCmd.Flags().IntP("limit", "l", 10, "Number of upcoming launches to show")
	pastCmd.Flags().IntP("limit", "l", 10, "Number of past launches to show")
}
