package controller

import (
	"encoding/json"
	"moderngo/repository"
	"net/http"
	"strings"
	"sync"
)

type WeatherController struct {
	repo *repository.WeatherRepository
}

func NewWeatherController(repo *repository.WeatherRepository) *WeatherController {
	return &WeatherController{repo: repo}
}

func (c *WeatherController) GetWeather(w http.ResponseWriter, r *http.Request) {
	lat := r.URL.Query().Get("lat")
	lon := r.URL.Query().Get("lon")

	if lat == "" || lon == "" {
		lat, lon = "52.52", "13.41" // Default to Berlin
	}

	data, err := c.repo.FetchWeather(lat, lon, "sample-token")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (c *WeatherController) BatchGet(w http.ResponseWriter, r *http.Request) {
	cities := strings.Split(r.URL.Query().Get("cities"), ",")
	if len(cities) == 0 || cities[0] == "" {
		cities = []string{"52.52:13.41", "48.85:2.35"} // Berlin, Paris
	}

	var wg sync.WaitGroup
	results := make(chan *repository.WeatherResponse, len(cities))

	for _, city := range cities {
		coords := strings.Split(city, ":")
		if len(coords) != 2 {
			continue
		}

		wg.Add(1)
		go func(lat, lon string) {
			defer wg.Done()
			if data, err := c.repo.FetchWeather(lat, lon, "batch-token"); err == nil {
				results <- data
			}
		}(coords[0], coords[1])
	}

	wg.Wait()
	close(results)

	var output []*repository.WeatherResponse
	for res := range results {
		output = append(output, res)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
