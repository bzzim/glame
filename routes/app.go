package routes

import (
	"github.com/bzzim/glame/controllers"
	"github.com/gin-gonic/gin"
)

type RouteAppController struct {
	controller          controllers.AppController
	authMiddleware      gin.HandlerFunc
	checkAuthMiddleware gin.HandlerFunc
}

// TODO: authMiddleware и checkAuthMiddleware
// И вообще их объеденить как-то
func NewRouteAppController(c controllers.AppController, am gin.HandlerFunc, cm gin.HandlerFunc) RouteAppController {
	return RouteAppController{
		controller:          c,
		authMiddleware:      am,
		checkAuthMiddleware: cm,
	}
}

func (r *RouteAppController) AddRoutes(router *gin.RouterGroup) {
	AppGroup := router.Group("apps")
	AppGroup.GET("", r.checkAuthMiddleware, r.controller.App)
}
