package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock HTTP server for testing fetchWeather function
func TestFetchWeather(t *testing.T) {
	// Mock API response
	mockResponse := `{
		"name": "London",
		"main": {
			"temp": 15.0,
			"humidity": 80
		},
		"wind": {
			"speed": 5.5
		}
	}`

	// Start a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Override baseURL for the test
	baseURL = server.URL

	// Call the function
	weather, err := fetchWeather("London")
	if err != nil {
		t.Fatalf("fetchWeather returned an error: %v", err)
	}

	// Validate the response
	if weather.Name != "London" {
		t.Errorf("expected city name London, got %s", weather.Name)
	}
	if weather.Main.Temp != 15.0 {
		t.Errorf("expected temperature 15.0, got %.1f", weather.Main.Temp)
	}
	if weather.Main.Humidity != 80 {
		t.Errorf("expected humidity 80, got %d", weather.Main.Humidity)
	}
	if weather.Wind.Speed != 5.5 {
		t.Errorf("expected wind speed 5.5, got %.1f", weather.Wind.Speed)
	}
}

// Test weatherHandler with an HTTP request
func TestWeatherHandler(t *testing.T) {
	// Mock API response
	mockResponse := `{
		"name": "London",
		"main": {
			"temp": 15.0,
			"humidity": 80
		},
		"wind": {
			"speed": 5.5
		}
	}`

	// Start a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Override baseURL for the test
	baseURL = server.URL

	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/weather?city=London", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	// Create a ResponseRecorder to capture the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(weatherHandler)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Validate the response
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var weather WeatherResponse
	err = json.NewDecoder(rr.Body).Decode(&weather)
	if err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if weather.Name != "London" {
		t.Errorf("expected city name London, got %s", weather.Name)
	}
	if weather.Main.Temp != 15.0 {
		t.Errorf("expected temperature 15.0, got %.1f", weather.Main.Temp)
	}
	if weather.Main.Humidity != 80 {
		t.Errorf("expected humidity 80, got %d", weather.Main.Humidity)
	}
	if weather.Wind.Speed != 5.5 {
		t.Errorf("expected wind speed 5.5, got %.1f", weather.Wind.Speed)
	}
}
