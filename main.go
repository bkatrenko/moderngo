package main

import (
	"fmt"
	"log"
	"moderngo/controller"
	"moderngo/repository"
	"net/http"
	"net/http/pprof"
	"time"
)

func main() {
	repo := repository.NewWeatherRepository()
	ctrl := controller.NewWeatherController(repo)

	http.HandleFunc("/weather", ctrl.GetWeather)
	http.HandleFunc("/weather/batch", ctrl.BatchGet)

	// Standard pprof endpoint: Target for Flight Recorder / goroutineleak additions
	http.HandleFunc("/debug/pprof/", pprof.Index)

	// Real-time polling loop: Target for testing with synctest
	go func() {
		for {
			time.Sleep(10 * time.Second)
		}
	}()

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
