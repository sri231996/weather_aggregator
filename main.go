package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	apiKey  = "your_openweathermap_api_key"
	baseURL = "https://api.openweathermap.org/data/2.5/weather"
)

type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

// Fetch weather data directly from the API
func fetchWeather(city string) (*WeatherResponse, error) {
	url := fmt.Sprintf("%s?q=%s&appid=%s&units=metric", baseURL, city, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch weather: %s", resp.Status)
	}

	var weather WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		return nil, err
	}
	return &weather, nil
}

// HTTP handler for weather requests
func weatherHandler(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "City is required", http.StatusBadRequest)
		return
	}

	weather, err := fetchWeather(city)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(weather)
}

func main() {
	apiKey = os.Getenv("OPENWEATHER_API_KEY") // Set API key via environment variable
	if apiKey == "" {
		log.Fatal("API key not found. Set the OPENWEATHER_API_KEY environment variable.")
	}

	http.HandleFunc("/weather", weatherHandler)
	port := "8080"
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
