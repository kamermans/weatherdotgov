package weatherdotgov

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
)

var (
	// Timeout sets the maximum time for a weather request
	Timeout = 10 * time.Second

	httpClient = http.Client{
		Timeout: Timeout,
	}
)

// WeatherData is a detailed, point-in-time weather report
type WeatherData struct {
	XMLName               xml.Name     `xml:"current_observation" json:"-"`
	Credit                string       `xml:"credit"`                  // ex: NOAA's National Weather Service
	CreditURL             string       `xml:"credit_URL"`              // ex: http://weather.gov/
	Image                 WeatherImage `xml:"image"`                   // ex: <image></image>
	SuggestedPickup       string       `xml:"suggested_pickup"`        // ex: 15 minutes after the hour
	SuggestedPickupPeriod int          `xml:"suggested_pickup_period"` // ex: 60
	Location              string       `xml:"location"`                // ex: Winchester Regional, VA
	StationID             string       `xml:"station_id"`              // ex: KOKV
	Latitude              float64      `xml:"latitude"`                // ex: 39.15
	Longitude             float64      `xml:"longitude"`               // ex: -78.15
	ObservationTimeString string       `xml:"observation_time"`        // ex: Last Updated on Oct 11 2020, 11:55 am EDT
	ObservationTimeRfc822 string       `xml:"observation_time_rfc822"` // ex: Sun, 11 Oct 2020 11:55:00 -0400
	ObservationTime       time.Time    // Go-native time, added after decoding
	Weather               string       `xml:"weather"`             // ex: Mostly Cloudy
	TemperatureString     string       `xml:"temperature_string"`  // ex: 70.0 F (21.0 C)
	TempF                 float64      `xml:"temp_f"`              // ex: 70.0
	TempC                 float64      `xml:"temp_c"`              // ex: 21.0
	RelativeHumidity      float64      `xml:"relative_humidity"`   // ex: 73
	WindString            string       `xml:"wind_string"`         // ex: Calm
	WindDir               string       `xml:"wind_dir"`            // ex: North
	WindDegrees           float64      `xml:"wind_degrees"`        // ex: 0
	WindMph               float64      `xml:"wind_mph"`            // ex: 0.0
	WindKt                float64      `xml:"wind_kt"`             // ex: 0
	PressureIn            float64      `xml:"pressure_in"`         // ex: 30.03
	DewpointString        string       `xml:"dewpoint_string"`     // ex: 60.8 F (16.0 C)
	DewpointF             float64      `xml:"dewpoint_f"`          // ex: 60.8
	DewpointC             float64      `xml:"dewpoint_c"`          // ex: 16.0
	VisibilityMi          float64      `xml:"visibility_mi"`       // ex: 10.00
	IconURLBase           string       `xml:"icon_url_base"`       // ex: http://forecast.weather.gov/images/wtf/small/
	TwoDayHistoryURL      string       `xml:"two_day_history_url"` // ex: http://www.weather.gov/data/obhistory/KOKV.html
	IconURLName           string       `xml:"icon_url_name"`       // ex: bkn.png
	ObURL                 string       `xml:"ob_url"`              // ex: http://www.weather.gov/data/METAR/KOKV.1.txt
	DisclaimerURL         string       `xml:"disclaimer_url"`      // ex: http://weather.gov/disclaimer.html
	CopyrightURL          string       `xml:"copyright_url"`       // ex: http://weather.gov/disclaimer.html
	PrivacyPolicyURL      string       `xml:"privacy_policy_url"`  // ex: http://weather.gov/notice.html
}

// WeatherImage contains an image for the weather stations
type WeatherImage struct {
	URL   string `xml:"url"`
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

func (w *WeatherData) String() string {
	out, _ := json.MarshalIndent(w, "", "  ")
	return string(out)
}

// CurrentWeatherFromStation gets the current weather for the given station.
// You can find your nearest station here: https://w1.weather.gov/xml/current_obs/seek.php
func CurrentWeatherFromStation(stationID string) (*WeatherData, error) {
	xmlURL := fmt.Sprintf("https://w1.weather.gov/xml/current_obs/%v.xml", stationID)
	return CurrentWeather(xmlURL)
}

// CurrentWeather gets the current weather for the given XML URL.
// You can find your nearest xml URL here: https://w1.weather.gov/xml/current_obs/seek.php
func CurrentWeather(xmlURL string) (*WeatherData, error) {

	resp, err := httpClient.Get(xmlURL)
	if err != nil {
		return nil, fmt.Errorf("unable to get weather data: %w", err)
	}

	data, err := parseWeatherData(resp.Body)

	// Make sure the entire body is read
	ioutil.ReadAll(resp.Body)

	return data, err
}

func parseWeatherData(reader io.Reader) (*WeatherData, error) {
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	data := new(WeatherData)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	// It says RFC822 in the name, but it's RFC1123Z when I check it
	data.ObservationTime, err = time.Parse(time.RFC1123Z, data.ObservationTimeRfc822)
	if err != nil {
		// Try RFC8222Z just in case ...
		data.ObservationTime, err = time.Parse(time.RFC822Z, data.ObservationTimeRfc822)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
