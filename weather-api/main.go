package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type WeatherResponse struct {
	Name string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
	} `json:"sys"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Main string `json:"main"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	city := "Dhaka"

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		fmt.Println("API key not found in environment")
		return
	}

	url := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s&units=metric",
		city,
		apiKey,
	)

	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("Network error:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Println("City not found (404)")
		return
	}

	if resp.StatusCode == http.StatusTooManyRequests {
		fmt.Println("Rate limit exceeded (429)")
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}

	var weather WeatherResponse
	err = json.Unmarshal(body, &weather)
	if err != nil {
		fmt.Println("Malformed JSON response:", err)
		return
	}

	if len(weather.Weather) == 0 {
		fmt.Println("Weather data unavailable")
		return
	}

	condition := weather.Weather[0].Main

	fmt.Println("Current Weather Report")
	fmt.Println("----------------------")
	fmt.Printf("City: %s, %s\n", weather.Name, weather.Sys.Country)
	fmt.Printf("Temperature: %.2f °C\n", weather.Main.Temp)
	fmt.Printf("Condition: %s\n", condition)
	fmt.Printf("Humidity: %d %%\n", weather.Main.Humidity)
	fmt.Printf("Wind Speed: %.2f m/s\n", weather.Wind.Speed)
}
