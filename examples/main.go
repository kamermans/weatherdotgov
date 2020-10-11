package main

import (
	"fmt"

	"github.com/kamermans/weatherdotgov"
)

func main() {
	w, err := weatherdotgov.CurrentWeatherFromStation("NSTU")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current Weather in %v\n", w.Location)
	fmt.Printf("Lat/Long: %v, %v\n", w.Latitude, w.Longitude)
	fmt.Printf("Temperature: %v F\n", w.TempF)
	fmt.Printf("Relative Humidity: %v %%\n", w.RelativeHumidity)
	fmt.Printf("Air Pressure: %v in\n", w.PressureIn)

	fmt.Printf("Complete report:\n%s\n", w)
}
