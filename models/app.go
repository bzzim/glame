package models

import "github.com/SamuelTissot/sqltime"

type App struct {
	Id          int          `json:"id" gorm:"column:id"`
	Name        string       `json:"name" gorm:"column:name" form:"name"`
	Url         string       `json:"url" gorm:"column:url" form:"url"`
	Icon        string       `json:"icon" gorm:"column:icon"`
	IsPublic    bool         `json:"isPublic" gorm:"column:isPublic" form:"isPublic"`
	IsPinned    bool         `json:"isPinned" gorm:"column:isPinned"`
	OrderId     int          `json:"orderId" gorm:"column:orderId"`
	Description string       `json:"description" gorm:"column:description" form:"description"`
	CreatedAt   sqltime.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt   sqltime.Time `json:"updatedAt" gorm:"column:updatedAt"`
}
