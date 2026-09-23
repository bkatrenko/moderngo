package main

import (
	"fmt"
	"log"
	"moderngo/controller"
	"moderngo/repository"
	"net/http"
)

func main() {
	repo := repository.NewWeatherRepository()
	ctrl := controller.NewWeatherController(repo)

	http.HandleFunc("/weather", ctrl.GetWeather)
	http.HandleFunc("/weather/batch", ctrl.BatchGet)

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
