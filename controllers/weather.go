package controllers

import (
	"fmt"
	"github.com/bzzim/glame/controllers/setting"
	"github.com/bzzim/glame/helper"
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/pkg/weather"
	"github.com/bzzim/glame/responses"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
	"reflect"
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
	r.db.Order("createdAt desc").Limit(1).Find(&rows)
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
		TempC:         data.TempC,
		TempF:         data.TempF,
		IsDay:         data.IsDay,
		Cloud:         data.Cloud,
		ConditionText: data.ConditionText,
		ConditionCode: data.ConditionCode,
		Humidity:      data.Humidity,
		WindK:         data.WindK,
		WindM:         data.WindM,
	}
	r.db.Create(&row)

	responses.NewSuccessResponse(ctx, row)
}

// TODO: вынести в хелпер
func (r *WeatherController) getFloat(unk interface{}) (float64, error) {
	var floatType = reflect.TypeOf(float64(0))
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}
	fv := v.Convert(floatType)
	return fv.Float(), nil
}
