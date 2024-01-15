package controllers

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/SamuelTissot/sqltime"
	"github.com/bzzim/glame/helper"
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/requests"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BookmarkController struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewBookmarkController(db *gorm.DB, log *slog.Logger) BookmarkController {
	return BookmarkController{db: db, log: log}
}

func (r *BookmarkController) AddBookmark(ctx *gin.Context) {
	var request requests.Bookmark

	file, err := ctx.FormFile("icon")
	// если передан файл, значит биндим запрос как форму и загружаем файл
	if err == nil {
		if err := ctx.ShouldBind(&request); err != nil {
			responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
			return
		}

		iconName, err := helper.UploadFile(file, ctx)
		if err != nil {
			responses.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		request.Icon = iconName
	} else {
		// а иначе биндим как json
		if err := ctx.ShouldBindJSON(&request); err != nil {
			responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
			return
		}
	}

	row := models.Bookmark{
		Name:       request.Name,
		Url:        request.Url,
		Icon:       strings.TrimSpace(request.Icon),
		CategoryId: request.CategoryId,
		IsPublic:   request.IsPublic,
		CreatedAt:  sqltime.Now(),
		UpdatedAt:  sqltime.Now(),
	}

	if result := r.db.Create(&row); result.Error != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}
	responses.NewSuccessCodeResponse(ctx, http.StatusCreated, row)
}

func (r *BookmarkController) SaveBoomark(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))

	if err != nil || id == 0 {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	var request models.Bookmark
	file, err := ctx.FormFile("icon")
	// если передан файл, значит биндим запрос как форму и загружаем файл
	if err == nil {
		if err := ctx.ShouldBind(&request); err != nil {
			responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
			return
		}

		iconName, err := helper.UploadFile(file, ctx)
		if err != nil {
			responses.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		request.Icon = iconName
	} else {
		// а иначе биндим как json
		if err := ctx.ShouldBindJSON(&request); err != nil {
			responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
			return
		}
	}

	row := models.Bookmark{Id: id}
	if err := r.db.First(&row).Error; err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	if len(request.Name) != 0 {
		row.Name = request.Name
	}
	if len(request.Url) != 0 {
		row.Url = request.Url
	}
	row.Icon = strings.TrimSpace(request.Icon)
	if request.CategoryId != 0 {
		row.CategoryId = request.CategoryId
	}
	row.IsPublic = request.IsPublic
	row.UpdatedAt = sqltime.Now()

	if err := r.db.Omit("createdAt").Save(&row).Error; err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	responses.NewSuccessResponse(ctx, row)
}

func (r *BookmarkController) DeleteBookmark(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	if err := r.db.Delete(models.Bookmark{Id: id}).Error; err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	responses.NewSuccessResponse(ctx, nil)
}

func (r *BookmarkController) ReorderBookmark(ctx *gin.Context) {
	var request requests.BookmarkOrder

	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	for _, b := range request.Bookmarks {
		r.db.Model(&models.Bookmark{Id: b.Id}).Updates(models.Bookmark{OrderId: b.OrderId})
	}

	responses.NewSuccessResponse(ctx, nil)
}
