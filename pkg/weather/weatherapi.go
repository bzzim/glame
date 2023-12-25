package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type WeatherApiResponse struct {
	Current struct {
		LastUpdated string  `json:"last_updated"`
		TempC       float32 `json:"temp_c"`
		TempF       float32 `json:"temp_f"`
		IsDay       int     `json:"is_day"`
		Condition   struct {
			Text string `json:"text"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMph  float32 `json:"wind_mph"`
		WindKph  float32 `json:"wind_kph"`
		Humidity int     `json:"humidity"`
		Cloud    int     `json:"cloud"`
	} `json:"current"`
}

type WeatherApi struct {
}

func (r *WeatherApi) Get(secret string, lat, lon float64) (*Weather, error) {
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%f,%f", secret, lat, lon)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data WeatherApiResponse
	body, err := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	result := Weather{
		ExtLastUpdate: data.Current.LastUpdated,
		TempC:         data.Current.TempC,
		TempF:         data.Current.TempF,
		IsDay:         data.Current.IsDay,
		Cloud:         data.Current.Cloud,
		ConditionText: data.Current.Condition.Text,
		ConditionCode: data.Current.Condition.Code,
		Humidity:      data.Current.Humidity,
		WindK:         data.Current.WindKph,
		WindM:         data.Current.WindMph,
	}
	return &result, nil
}
