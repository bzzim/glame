package controllers

import (
	"github.com/SamuelTissot/sqltime"
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/requests"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

type BookmarkController struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewBookmarkController(db *gorm.DB, log *slog.Logger) BookmarkController {
	return BookmarkController{db: db, log: log}
}

func (r *CategoryController) AddBookmark(ctx *gin.Context) {
	var request requests.Category
	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	row := models.Category{
		Name:      request.Name,
		IsPublic:  request.IsPublic,
		CreatedAt: sqltime.Now(),
		UpdatedAt: sqltime.Now(),
	}
	if result := r.db.Create(&row); result.Error != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}
	responses.NewSuccessResponse(ctx, row)
}
