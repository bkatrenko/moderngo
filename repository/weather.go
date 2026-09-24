package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

type Priority int

const (
	defaultPriority Priority = iota + 1
	special
	biggestPriority
)

type CurrentWeather struct {
	Temperature float64 `json:"temperature"`
	WindSpeed   float64 `json:"windspeed"`
	Time        string  `json:"time"`
}

type WeatherResponse struct {
	Latitude       float64        `json:"latitude"`
	Longitude      float64        `json:"longitude"`
	CurrentWeather CurrentWeather `json:"current_weather"`
	CustomPriority *Priority      `json:"custom_priority,omitempty"`
	ErrorMessage   string         `json:"error_message,omitempty"`
}

type WeatherRepository struct {
	mu sync.RWMutex

	cache map[string]*WeatherResponse
}

func NewWeatherRepository() *WeatherRepository {
	return &WeatherRepository{
		cache: make(map[string]*WeatherResponse),
	}
}

func (r *WeatherRepository) FetchWeather(ctx context.Context, lat, lon string) (*WeatherResponse, error) {
	cacheKey := fmt.Sprintf("%s:%s", lat, lon)

	r.mu.RLock()
	if cached, ok := r.cache[cacheKey]; ok {
		r.mu.RUnlock()
		return cached, nil
	}
	r.mu.RUnlock()

	url := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current_weather=true", lat, lon)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weather WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, err
	}

	weather.CustomPriority = new(Priority(defaultPriority))

	r.mu.Lock()
	r.cache[cacheKey] = &weather
	r.mu.Unlock()

	return &weather, nil
}

func (r *WeatherRepository) GetAllCached() []WeatherResponse {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]WeatherResponse, 0, len(r.cache))
	for _, v := range r.cache {
		results = append(results, *v)
	}
	return results
}
