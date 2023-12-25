package controllers

import (
	"github.com/SamuelTissot/sqltime"
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"time"
)

type AppController struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewAppController(db *gorm.DB, log *slog.Logger) AppController {
	return AppController{db: db, log: log}
}

func (r *AppController) App(ctx *gin.Context) {
	sqltime.TruncateOff = time.Microsecond
	//isAuth := ctx.GetBool("isAuth")
	var bookmarks []models.Bookmark
	var query interface{}
	//if !isAuth {
	//	query = &models.App{IsPublic: 1}
	//}
	// TODO:  сортировка
	r.db.Where(query).Find(&bookmarks)
	responses.NewSuccessResponse(ctx, bookmarks)
}
