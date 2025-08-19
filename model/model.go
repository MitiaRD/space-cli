package model

import "time"

type Launch struct {
	ID           string    `json:"id"`
	FlightNumber int       `json:"flight_number"`
	Name         string    `json:"name"`
	Date         time.Time `json:"date_utc"`
	Success      *bool     `json:"success"`
	Crew         []string  `json:"crew"`
	RocketId     string    `json:"rocket"`
	Details      string    `json:"details"`
	LaunchpadId  string    `json:"launchpad"`
}

type Launchpad struct {
	ID        string  `json:"id"`
	Name      string  `json:"full_name"`
	Locality  string  `json:"locality"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
	Details   string  `json:"details"`
}

type Crew struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Agency string `json:"agency"`
	Image  string `json:"image"`
	Status string `json:"status"`
}

type Rocket struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CostPerLaunch int    `json:"cost_per_launch"`
	SuccessRate   int    `json:"success_rate_pct"`
	Country       string `json:"country"`
	Company       string `json:"company"`
	Height        Length `json:"height_w_trunk"`
	Diameter      Length `json:"diameter"`
	Mass          Mass   `json:"mass"`
	FirstFlight   string `json:"first_flight"`
	Description   string `json:"description"`
}

type Mass struct {
	Kg float32 `json:"kg"`
	Lb float32 `json:"lb"`
}

type Length struct {
	Meters float32 `json:"meters"`
	Feet   float32 `json:"feet"`
}

// NASA API Models
type NasaEarthEvent struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Categories  []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"categories"`
	Geometry []struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
		Date        string    `json:"date"`
	} `json:"geometry"`
}

type NasaEarth struct {
	Events []NasaEarthEvent `json:"events"`
}
