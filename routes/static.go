package routes

import (
	"github.com/gin-gonic/gin"
)

type RouteStaticController struct {
}

func NewRouteStaticController() RouteStaticController {
	return RouteStaticController{}
}

func (r *RouteStaticController) AddRoutes(router *gin.Engine) {
	router.Static("uploads", "./data/uploads")
	router.StaticFile("flame.css", "./data/flame.css")
}
