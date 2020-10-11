# Weather.gov Current Weather Observations

This simple Go library retrieves and parses XML (not RSS) weather observations from weather.gov (the US National Weather Service).

## Usage
First, find your nearest station by visiting the [National Weather Service Station Finder](https://w1.weather.gov/xml/current_obs/seek.php) and selecting your US state under **XML Feeds of Current Weather Conditions**.

For example, if you pick **American Samoa**, you will see the station `Pago Pago, AS, Samoa` with its station ID `NSTU`.

```go
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
}
```

Output:
```
Current Weather in Pago Pago, AS, Samoa, AS
Lat/Long: -14.331, -170.7105
Temperature: 82 F
Relative Humidity: 81 %
Air Pressure: 29.88 in
```

## JSON Format
You can also get the weather report in JSON format by using the `WeatherData`'s `String()` function:
```go
fmt.Println(w.String())
```

```json
{
  "Credit": "NOAA's National Weather Service",
  "CreditURL": "http://weather.gov/",
  "Image": {
    "URL": "http://weather.gov/images/xml_logo.gif",
    "Title": "NOAA's National Weather Service",
    "Link": "http://weather.gov"
  },
  "SuggestedPickup": "15 minutes after the hour",
  "SuggestedPickupPeriod": 60,
  "Location": "Pago Pago, AS, Samoa, AS",
  "StationID": "NSTU",
  "Latitude": -14.331,
  "Longitude": -170.7105,
  "ObservationTimeString": "Last Updated on Oct 11 2020, 6:50 am SST",
  "ObservationTimeRfc822": "Sun, 11 Oct 2020 06:50:00 -1100",
  "ObservationTime": "2020-10-11T06:50:00-11:00",
  "Weather": "Mostly Cloudy",
  "TemperatureString": "82.0 F (27.7 C)",
  "TempF": 82,
  "TempC": 27.7,
  "RelativeHumidity": 81,
  "WindString": "East at 16.1 MPH (14 KT)",
  "WindDir": "East",
  "WindDegrees": 110,
  "WindMph": 16.1,
  "WindKt": 14,
  "PressureIn": 29.88,
  "DewpointString": "75.6 F (24.2 C)",
  "DewpointF": 75.6,
  "DewpointC": 24.2,
  "VisibilityMi": 12,
  "IconURLBase": "http://forecast.weather.gov/images/wtf/small/",
  "TwoDayHistoryURL": "http://www.weather.gov/data/obhistory/NSTU.html",
  "IconURLName": "nbkn.png",
  "ObURL": "http://www.weather.gov/data/METAR/NSTU.1.txt",
  "DisclaimerURL": "http://weather.gov/disclaimer.html",
  "CopyrightURL": "http://weather.gov/disclaimer.html",
  "PrivacyPolicyURL": "http://weather.gov/notice.html"
}
```