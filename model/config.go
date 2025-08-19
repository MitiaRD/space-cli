package model

import (
	"fmt"
	"time"
)

type Config struct {
	NASAAPIKey string        `validate:"required"`
	Timeout    time.Duration `validate:"required,min=1s"`
	Retries    int           `validate:"required,min=1,max=10"`
	BaseDelay  time.Duration `validate:"required,min=100ms"`
	MaxDelay   time.Duration `validate:"required,min=1s"`
}

func (c *Config) Validate() error {
	if c.NASAAPIKey == "" {
		return fmt.Errorf("NASA API key is required")
	}
	if c.Timeout < time.Second {
		return fmt.Errorf("timeout must be at least 1 second")
	}
	if c.Retries < 1 || c.Retries > 10 {
		return fmt.Errorf("retries must be between 1 and 10")
	}
	if c.BaseDelay < 100*time.Millisecond {
		return fmt.Errorf("base delay must be at least 100ms")
	}
	if c.MaxDelay < time.Second {
		return fmt.Errorf("max delay must be at least 1 second")
	}
	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Timeout:   30 * time.Second,
		Retries:   5,
		BaseDelay: 200 * time.Millisecond,
		MaxDelay:  5 * time.Second,
	}
}
