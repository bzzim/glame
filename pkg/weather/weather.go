package weather

import "time"

type Weather struct {
	Id            int       `json:"id"`
	ExtLastUpdate string    `json:"externalLastUpdate"`
	TempC         float32   `json:"tempC"`
	TempF         float32   `json:"tempF"`
	IsDay         int       `json:"isDay"`
	Cloud         int       `json:"cloud"`
	ConditionText string    `json:"conditionText"`
	ConditionCode int       `json:"conditionCode"`
	Humidity      int       `json:"humidity"`
	WindK         float32   `json:"windK"`
	WindM         float32   `json:"windM"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Service interface {
	Get(secret string, lat, lon float64) (*Weather, error)
}
