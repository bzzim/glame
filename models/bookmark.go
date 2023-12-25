package models

import "time"

type Bookmark struct {
	Id         int       `json:"id" gorm:"column:id"`
	Name       string    `json:"name" gorm:"column:name"`
	Url        string    `json:"url" gorm:"column:url"`
	CategoryId int       `json:"-" gorm:"column:category_id"`
	Icon       string    `json:"icon" gorm:"column:icon"`
	CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at"`
	IsPublic   bool      `json:"isPublic" gorm:"column:is_public"`
	OrderId    int       `json:"orderId" gorm:"column:order_id"`
}
