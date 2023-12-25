package controllers

import (
	"fmt"
	"github.com/SamuelTissot/sqltime"
	"github.com/bzzim/glame/controllers/setting"
	"github.com/bzzim/glame/helper"
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/requests"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"strconv"
)

type CategoryController struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewCategoryController(db *gorm.DB, log *slog.Logger) CategoryController {
	return CategoryController{db: db, log: log}
}

func (r *CategoryController) Categories(ctx *gin.Context) {
	//sqltime.TruncateOff = time.Microsecond
	isAuth := ctx.GetBool("isAuth")
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
	// TODO: в оригинале есл строит сортировка по имени, то еще и сами ссылки сортируем по имени
	orderField := "name"
	if cfg.UseOrdering == "createdAt" {
		orderField = "created_at"
	} else if cfg.UseOrdering == "orderId" {
		orderField = "order_id"
	}

	order := fmt.Sprintf("%s asc", orderField)
	r.db.Order(order).Where(query).Preload("Bookmarks").Find(&categories)
	responses.NewSuccessResponse(ctx, categories)
}

func (r *CategoryController) AddCategory(ctx *gin.Context) {
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

	if err := r.db.Save(&row).Error; err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	responses.NewSuccessResponse(ctx, row)
}

func (r *CategoryController) OrderCategory(ctx *gin.Context) {
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
