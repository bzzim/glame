package setting

import (
	"gorm.io/gorm"
	"log/slog"
)

type Controller struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewController(db *gorm.DB, log *slog.Logger) Controller {
	return Controller{db: db, log: log}
}
