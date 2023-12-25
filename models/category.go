package models

import (
	"github.com/SamuelTissot/sqltime"
)

type Category struct {
	Id        int          `json:"id" gorm:"primarykey,column:id"`
	Name      string       `json:"name" gorm:"column:name"`
	IsPublic  bool         `json:"isPublic" gorm:"column:is_public"`
	IsPinned  bool         `json:"isPinned" gorm:"column:is_pinned"`
	OrderId   int          `json:"orderId" gorm:"column:order_id"`
	CreatedAt sqltime.Time `json:"createdAt" gorm:"column:created_at;type:timestamp"`
	UpdatedAt sqltime.Time `json:"updatedAt" gorm:"column:updated_at;type:timestamp"`
	Bookmarks []Bookmark   `json:"bookmarks" gorm:"foreignKey:category_id"`
}
