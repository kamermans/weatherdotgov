package weatherdotgov

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testSample = `
<?xml version="1.0" encoding="ISO-8859-1"?>
<?xml-stylesheet href="latest_ob.xsl" type="text/xsl"?>
<current_observation version="1.0"
         xmlns:xsd="http://www.w3.org/2001/XMLSchema"
         xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:noNamespaceSchemaLocation="http://www.weather.gov/view/current_observation.xsd">
        <credit>NOAA's National Weather Service</credit>
        <credit_URL>http://weather.gov/</credit_URL>
        <image>
                <url>http://weather.gov/images/xml_logo.gif</url>
                <title>NOAA's National Weather Service</title>
                <link>http://weather.gov</link>
        </image>
        <suggested_pickup>15 minutes after the hour</suggested_pickup>
        <suggested_pickup_period>60</suggested_pickup_period>
        <location>Winchester Regional, VA</location>
        <station_id>KOKV</station_id>
        <latitude>39.15</latitude>
        <longitude>-78.15</longitude>
        <observation_time>Last Updated on Oct 11 2020, 11:55 am EDT</observation_time>
        <observation_time_rfc822>Sun, 11 Oct 2020 11:55:00 -0400</observation_time_rfc822>
        <weather>Mostly Cloudy</weather>
        <temperature_string>70.0 F (21.0 C)</temperature_string>
        <temp_f>70.0</temp_f>
        <temp_c>21.0</temp_c>
        <relative_humidity>73</relative_humidity>
        <wind_string>Calm</wind_string>
        <wind_dir>North</wind_dir>
        <wind_degrees>0</wind_degrees>
        <wind_mph>0.0</wind_mph>
        <wind_kt>0</wind_kt>
        <pressure_in>30.03</pressure_in>
        <dewpoint_string>60.8 F (16.0 C)</dewpoint_string>
        <dewpoint_f>60.8</dewpoint_f>
        <dewpoint_c>16.0</dewpoint_c>
        <visibility_mi>10.00</visibility_mi>
        <icon_url_base>http://forecast.weather.gov/images/wtf/small/</icon_url_base>
        <two_day_history_url>http://www.weather.gov/data/obhistory/KOKV.html</two_day_history_url>
        <icon_url_name>bkn.png</icon_url_name>
        <ob_url>http://www.weather.gov/data/METAR/KOKV.1.txt</ob_url>
        <disclaimer_url>http://weather.gov/disclaimer.html</disclaimer_url>
        <copyright_url>http://weather.gov/disclaimer.html</copyright_url>
        <privacy_policy_url>http://weather.gov/notice.html</privacy_policy_url>
</current_observation>
`

func TestParseWeatherData(t *testing.T) {
	reader := bytes.NewReader([]byte(testSample))
	w, err := parseWeatherData(reader)
	require.Nil(t, err)
	require.NotNil(t, w)

	assert.Equal(t, 73.0, w.RelativeHumidity)
	assert.Equal(t, 60.8, w.DewpointF)
	assert.Equal(t, "Sun, 11 Oct 2020 11:55:00 -0400", w.ObservationTimeRfc822)

	expected, err := time.Parse(time.RFC1123, "Sun, 11 Oct 2020 11:55:00 EDT")
	require.Nil(t, err)
	assert.Equal(t, expected, w.ObservationTime)
}

func TestCurrentWeather(t *testing.T) {
	w, err := CurrentWeather("https://w1.weather.gov/xml/current_obs/PAMR.xml")
	require.Nil(t, err)
	require.NotNil(t, w)

	assert.Equal(t, "PAMR", w.StationID)
	assert.Equal(t, "Anchorage, Merrill Field Airport, AK", w.Location)
	assert.NotEqual(t, 0, w.RelativeHumidity)
	assert.NotEqual(t, 0, w.DewpointF)
}

func TestCurrentWeatherFromStation(t *testing.T) {
	w, err := CurrentWeatherFromStation("PAMR")
	require.Nil(t, err)
	require.NotNil(t, w)

	assert.Equal(t, "PAMR", w.StationID)
	assert.Equal(t, "Anchorage, Merrill Field Airport, AK", w.Location)
	assert.NotEqual(t, 0, w.RelativeHumidity)
	assert.NotEqual(t, 0, w.DewpointF)
}
