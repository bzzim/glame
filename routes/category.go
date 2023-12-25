package routes

import (
	"github.com/bzzim/glame/controllers"
	"github.com/gin-gonic/gin"
)

type RouteCategoryController struct {
	controller          controllers.CategoryController
	authMiddleware      gin.HandlerFunc
	checkAuthMiddleware gin.HandlerFunc
}

// TODO: authMiddleware и checkAuthMiddleware
// И вообще их объеденить как-то
func NewRouteCategoryController(c controllers.CategoryController, am gin.HandlerFunc, cm gin.HandlerFunc) RouteCategoryController {
	return RouteCategoryController{
		controller:          c,
		authMiddleware:      am,
		checkAuthMiddleware: cm,
	}
}

func (r *RouteCategoryController) AddRoutes(router *gin.RouterGroup) {
	categoryGroup := router.Group("categories")
	categoryGroup.GET("", r.checkAuthMiddleware, r.controller.Categories)
	categoryGroup.POST("", r.authMiddleware, r.controller.AddCategory)
	categoryGroup.PUT(":id", r.authMiddleware, r.controller.SaveCategory)
	categoryGroup.PUT("0/reorder", r.authMiddleware, r.controller.OrderCategory)
	categoryGroup.DELETE(":id", r.authMiddleware, r.controller.DeleteCategory)
}
