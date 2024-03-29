package controllers

import (
	"log/slog"
	"net/http"

	"github.com/bzzim/glame/controllers/setting"
	"github.com/bzzim/glame/helper"
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/pkg/utils"
	"github.com/bzzim/glame/pkg/weather"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WeatherController struct {
	db  *gorm.DB
	ws  weather.Service
	log *slog.Logger
}

func NewWeatherController(db *gorm.DB, ws weather.Service, log *slog.Logger) WeatherController {
	return WeatherController{db: db, ws: ws, log: log}
}

func (r *WeatherController) Weather(ctx *gin.Context) {
	var rows []models.Weather
	r.db.Table("weather").Order("createdAt desc").Limit(1).Find(&rows)
	responses.NewSuccessResponse(ctx, rows)
}

func (r *WeatherController) UpdateWeather(ctx *gin.Context) {
	cfg, err := helper.LoadConfig(setting.ConfigFile)
	if err != nil {
		r.log.Warn(err.Error())
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	data, err := r.ws.Get(cfg.WeatherAPIKey, cfg.Lat, cfg.Lon)
	if err != nil {
		r.log.Warn(err.Error())
		responses.NewErrorResponse(ctx, http.StatusInternalServerError, responses.ServerErrorMessage)
		return
	}

	row := models.Weather{
		ExtLastUpdate: data.ExtLastUpdate,
		TempC:         utils.Float64ToFixed(data.TempC, 1),
		TempF:         utils.Float64ToFixed(data.TempF, 1),
		IsDay:         data.IsDay,
		Cloud:         data.Cloud,
		ConditionText: data.ConditionText,
		ConditionCode: data.ConditionCode,
		Humidity:      data.Humidity,
		WindK:         utils.Float64ToFixed(data.WindK, 1),
		WindM:         utils.Float64ToFixed(data.WindM, 1),
	}
	r.db.Table("weather").Create(&row)

	responses.NewSuccessResponse(ctx, row)
}
