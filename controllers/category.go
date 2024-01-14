package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/SamuelTissot/sqltime"
	"github.com/bzzim/glame/controllers/setting"
	"github.com/bzzim/glame/helper"
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/requests"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryController struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewCategoryController(db *gorm.DB, log *slog.Logger) CategoryController {
	return CategoryController{db: db, log: log}
}

func (r *CategoryController) Categories(ctx *gin.Context) {
	isAuth := helper.UserIsAuth(ctx)
	var categories []models.Category
	var query interface{}
	if !isAuth {
		query = &models.Category{IsPublic: true}
	}

	cfg, err := helper.LoadConfig(setting.ConfigFile)
	if err != nil {
		r.log.Warn(err.Error())
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	orderField := "name"
	orderType := "asc"
	if cfg.UseOrdering == "createdAt" {
		orderField = "id"
		orderType = "desc"
	} else if cfg.UseOrdering == "orderId" {
		orderField = "orderId"
	}

	order := fmt.Sprintf("%s %s", orderField, orderType)
	r.db.Order(order).Where(query).Preload("Bookmarks", func(db *gorm.DB) *gorm.DB {
		db = db.Order(order)
		return db
	}).Find(&categories)
	responses.NewSuccessResponse(ctx, categories)
}

func (r *CategoryController) AddCategory(ctx *gin.Context) {
	cfg, err := helper.LoadConfig(setting.ConfigFile)
	if err != nil {
		r.log.Warn(err.Error())
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	var request requests.Category
	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	row := models.Category{
		Name:      request.Name,
		IsPublic:  request.IsPublic,
		IsPinned:  cfg.PinCategoriesByDefault,
		CreatedAt: sqltime.Now(),
		UpdatedAt: sqltime.Now(),
	}
	if result := r.db.Create(&row); result.Error != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}
	responses.NewSuccessResponse(ctx, row)
}

func (r *CategoryController) SaveCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	var request models.Category
	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	// TODO: к нижней тудушке - а может потому что тут всегда обновляется дата?
	row := models.Category{Id: id, UpdatedAt: sqltime.Now()}
	if err := r.db.First(&row).Error; err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	// TODO: на клиенте есть баг: после пересортировки не обновляется orderId,
	// если после этого скрыть/показать категорию, то он отправляет старый orderId
	// и соотвественно оно перезаписывается
	if request.Id == 0 {
		row.IsPinned = request.IsPinned
	} else {
		row = request
	}

	if err := r.db.Omit("createdAt").Save(&row).Error; err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	responses.NewSuccessResponse(ctx, row)
}

func (r *CategoryController) ReorderCategory(ctx *gin.Context) {
	var request struct {
		Categories []models.Category `json:"categories"`
	}
	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	for _, category := range request.Categories {
		r.db.Model(&category).Updates(models.Category{OrderId: category.OrderId})
	}

	responses.NewSuccessResponse(ctx, nil)
}

func (r *CategoryController) DeleteCategory(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	if err := r.db.Delete(models.Category{Id: id}).Error; err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	responses.NewSuccessResponse(ctx, nil)
}
