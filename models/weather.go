package models

import "github.com/SamuelTissot/sqltime"

type Weather struct {
	Id            uint         `json:"id" gorm:"column:id"`
	ExtLastUpdate string       `json:"externalLastUpdate" gorm:"column:externalLastUpdate"`
	TempC         float64      `json:"tempC" gorm:"column:tempC"`
	TempF         float64      `json:"tempF" gorm:"column:tempF"`
	IsDay         int          `json:"isDay" gorm:"column:isDay"`
	Cloud         int          `json:"cloud" gorm:"column:cloud"`
	ConditionText string       `json:"conditionText" gorm:"column:conditionText"`
	ConditionCode int          `json:"conditionCode" gorm:"column:conditionCode"`
	Humidity      int          `json:"humidity" gorm:"column:humidity"`
	WindK         float64      `json:"windK" gorm:"column:windK"`
	WindM         float64      `json:"windM" gorm:"column:windM"`
	CreatedAt     sqltime.Time `json:"createdAt" gorm:"column:createdAt;type:datetime"`
	UpdatedAt     sqltime.Time `json:"updatedAt" gorm:"column:updatedAt;type:datetime"`
}
