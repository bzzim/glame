package routes

import (
	"github.com/bzzim/glame/controllers"
	"github.com/gin-gonic/gin"
)

type WeatherRouteController struct {
	controller controllers.WeatherController
}

func NewWeatherRouteController(weatherController controllers.WeatherController) WeatherRouteController {
	return WeatherRouteController{weatherController}
}

func (r *WeatherRouteController) AddRoutes(router *gin.RouterGroup) {
	weatherGroup := router.Group("weather")
	weatherGroup.GET("", r.controller.Weather)
	weatherGroup.GET("update", r.controller.UpdateWeather)
}
