package routes

import (
	"github.com/bzzim/glame/controllers/setting"
	"github.com/gin-gonic/gin"
)

type RouteSettingController struct {
	controller          setting.Controller
	authMiddleware      gin.HandlerFunc
	checkAuthMiddleware gin.HandlerFunc
}

// RouteSettingController
// TODO: authMiddleware и checkAuthMiddleware
// И вообще их объеденить как-то
func NewRouteSettingController(c setting.Controller, am gin.HandlerFunc, cm gin.HandlerFunc) RouteSettingController {
	return RouteSettingController{
		controller:          c,
		authMiddleware:      am,
		checkAuthMiddleware: cm,
	}
}

func (r *RouteSettingController) AddRoutes(router *gin.RouterGroup) {
	configGroup := router.Group("config")
	configGroup.GET("", r.controller.Config)
	configGroup.GET("0/css", r.controller.Css)
	configGroup.PUT("", r.authMiddleware, r.controller.SaveConfig)
	configGroup.PUT("0/css", r.controller.SaveCss)

	themesGroup := router.Group("themes")
	themesGroup.GET("", r.controller.Themes)
	themesGroup.POST("", r.authMiddleware, r.controller.AddTheme)
	themesGroup.PUT(":name", r.authMiddleware, r.controller.SaveTheme)
	themesGroup.DELETE(":name", r.authMiddleware, r.controller.DeleteTheme)

	queriesGroup := router.Group("queries")
	queriesGroup.GET("", r.controller.Queries)
	queriesGroup.POST("", r.authMiddleware, r.controller.AddQuery)
	queriesGroup.PUT(":prefix", r.authMiddleware, r.controller.SaveQuery)
	queriesGroup.DELETE(":prefix", r.authMiddleware, r.controller.DeleteQuery)
}
