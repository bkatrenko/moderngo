package controller

import (
	"context"
	"encoding/json"
	"log"
	"moderngo/repository"
	"net/http"
	"strings"
	"sync"
	"time"
)

type WeatherRepository interface {
	FetchWeather(ctx context.Context, lat, lon string) (*repository.WeatherResponse, error)
}

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

	data, err := c.repo.FetchWeather(r.Context(), lat, lon)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func (c *WeatherController) BatchGet(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second*20)
	defer cancel()

	cities := strings.Split(r.URL.Query().Get("cities"), ",")
	if len(cities) == 0 || cities[0] == "" {
		cities = []string{"52.52:13.41", "48.85:2.35"} // Berlin, Paris
	}

	output, err := c.getCitiesWeatherInBatch(ctx, cities)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

func (c *WeatherController) getCitiesWeatherInBatch(ctx context.Context, cities []string) ([]*repository.WeatherResponse, error) {
	var wg sync.WaitGroup

	results := make(chan *repository.WeatherResponse, len(cities))

	for _, city := range cities {
		coords := strings.Split(city, ":")
		if len(coords) != 2 {
			continue
		}

		wg.Go(func() {
			data, err := c.repo.FetchWeather(ctx, coords[0], coords[1])
			if err != nil {
				log.Printf("error: %s\n", err)
				if data == nil {
					data = &repository.WeatherResponse{ErrorMessage: "unexpected error"}
				}
			}

			results <- data
		})
	}

	wg.Wait()
	close(results)

	var output []*repository.WeatherResponse
	for res := range results {
		output = append(output, res)
	}

	return output, nil
}
