/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/MitiaRD/ReMarkable-cli/api"
	"github.com/spf13/cobra"
)

// earthCmd represents the earth command
var earthCmd = &cobra.Command{
	Use:   "earth",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		events, err := api.GetEarthEvents("")
		if err != nil {
			fmt.Printf("Error fetching earth events: %v\n", err)
			return
		}
		fmt.Printf("Earth events: %+v\n", events)
	},
}

func init() {
	rootCmd.AddCommand(earthCmd)

	// Here you will define your flags and configuration settings.
	earthCmd.Flags().IntP("limit", "l", 10, "Number of past launches to show")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// earthCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// earthCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
