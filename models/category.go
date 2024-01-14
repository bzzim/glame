package models

import (
	"github.com/SamuelTissot/sqltime"
)

type Category struct {
	Id        int          `json:"id" gorm:"primarykey,column:id"`
	Name      string       `json:"name" gorm:"column:name"`
	IsPublic  bool         `json:"isPublic" gorm:"column:isPublic"`
	IsPinned  bool         `json:"isPinned" gorm:"column:isPinned"`
	OrderId   int          `json:"orderId" gorm:"column:orderId"`
	CreatedAt sqltime.Time `json:"createdAt" gorm:"column:createdAt;type:timestamp"`
	UpdatedAt sqltime.Time `json:"updatedAt" gorm:"column:updatedAt;type:timestamp"`
	Bookmarks []Bookmark   `json:"bookmarks" gorm:"foreignKey:categoryId"`
}
