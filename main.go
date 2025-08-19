/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MitiaRD/ReMarkable-cli/cmd"
)

func main() {
	logger := cmd.SetupLogger()
	logger.Info("Starting Space CLI application")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		logger.Info("Received shutdown signal", "signal", sig)
		os.Exit(0)
	}()

	cmd.Execute()
}
