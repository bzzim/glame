package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/SamuelTissot/sqltime"
	"github.com/bzzim/glame/controllers/setting"
	"github.com/bzzim/glame/helper"
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/requests"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AppController struct {
	db  *gorm.DB
	log *slog.Logger
}

func NewAppController(db *gorm.DB, log *slog.Logger) AppController {
	return AppController{db: db, log: log}
}

func (r *AppController) Apps(ctx *gin.Context) {
	isAuth := helper.UserIsAuth(ctx)
	var apps []models.App
	var query interface{}
	if !isAuth {
		query = &models.App{IsPublic: true}
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
	r.db.Order(order).Where(query).Find(&apps)
	responses.NewSuccessResponse(ctx, apps)
}

func (r *AppController) AddApp(ctx *gin.Context) {
	cfg, err := helper.LoadConfig(setting.ConfigFile)
	if err != nil {
		r.log.Warn(err.Error())
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	var request requests.App
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

	row := models.App{
		Name:        request.Name,
		Url:         request.Url,
		Icon:        strings.TrimSpace(request.Icon),
		IsPublic:    request.IsPublic,
		IsPinned:    cfg.PinAppsByDefault,
		Description: request.Description,
		CreatedAt:   sqltime.Now(),
		UpdatedAt:   sqltime.Now(),
	}

	if result := r.db.Create(&row); result.Error != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}
	responses.NewSuccessCodeResponse(ctx, http.StatusCreated, row)
}

func (r *AppController) SaveApp(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	var request models.App
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
		request.Id = id
		request.Icon = iconName
	} else {
		// а иначе биндим как json
		if err := ctx.ShouldBindJSON(&request); err != nil {
			responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
			return
		}
	}

	row := models.App{Id: id}
	if err := r.db.First(&row).Error; err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	if request.Id == 0 {
		row.IsPinned = request.IsPinned
	} else {
		row = request
	}
	row.UpdatedAt = sqltime.Now()

	if err := r.db.Omit("createdAt").Save(&row).Error; err != nil {
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	responses.NewSuccessResponse(ctx, row)
}

func (r *AppController) DeleteApp(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil || id == 0 {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	if err := r.db.Delete(models.App{Id: id}).Error; err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	responses.NewSuccessResponse(ctx, nil)
}

func (r *AppController) ReorderApp(ctx *gin.Context) {
	var request requests.AppOrder

	if err := ctx.ShouldBindJSON(&request); err != nil {
		responses.NewErrorResponse(ctx, http.StatusBadRequest, responses.BadRequestMessage)
		return
	}

	for _, b := range request.Apps {
		r.db.Model(&models.App{Id: b.Id}).Updates(models.App{OrderId: b.OrderId})
	}

	responses.NewSuccessResponse(ctx, nil)
}
