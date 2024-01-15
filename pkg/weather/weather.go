package weather

import "time"

type Weather struct {
	Id            int       `json:"id"`
	ExtLastUpdate string    `json:"externalLastUpdate"`
	TempC         float64   `json:"tempC"`
	TempF         float64   `json:"tempF"`
	IsDay         int       `json:"isDay"`
	Cloud         int       `json:"cloud"`
	ConditionText string    `json:"conditionText"`
	ConditionCode int       `json:"conditionCode"`
	Humidity      int       `json:"humidity"`
	WindK         float64   `json:"windK"`
	WindM         float64   `json:"windM"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Service interface {
	Get(secret string, lat, lon float64) (*Weather, error)
}
