package models

import "github.com/SamuelTissot/sqltime"

type Bookmark struct {
	Id         int          `json:"id" gorm:"column:id"`
	Name       string       `json:"name" gorm:"column:name"`
	Url        string       `json:"url" gorm:"column:url"`
	CategoryId int          `json:"categoryId" gorm:"column:categoryId"`
	Icon       string       `json:"icon" gorm:"column:icon"`
	IsPublic   bool         `json:"isPublic" gorm:"column:isPublic"`
	OrderId    int          `json:"orderId" gorm:"column:orderId"`
	CreatedAt  sqltime.Time `json:"createdAt" gorm:"column:createdAt"`
	UpdatedAt  sqltime.Time `json:"updatedAt" gorm:"column:updatedAt"`
}
