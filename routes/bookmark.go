package routes

import (
	"github.com/bzzim/glame/controllers"
	"github.com/gin-gonic/gin"
)

type RouteBookmarkController struct {
	controller          controllers.BookmarkController
	authMiddleware      gin.HandlerFunc
	checkAuthMiddleware gin.HandlerFunc
}

func NewRouteBookmarkController(c controllers.BookmarkController, am gin.HandlerFunc, cm gin.HandlerFunc) RouteBookmarkController {
	return RouteBookmarkController{
		controller:          c,
		authMiddleware:      am,
		checkAuthMiddleware: cm,
	}
}

func (r *RouteBookmarkController) AddRoutes(router *gin.RouterGroup) {
	bookmarkGroup := router.Group("bookmarks")
	bookmarkGroup.POST("", r.authMiddleware, r.controller.AddBookmark)
	bookmarkGroup.PUT(":id", r.authMiddleware, r.controller.SaveBoomark)
	bookmarkGroup.PUT("0/reorder", r.authMiddleware, r.controller.ReorderBookmark)
	bookmarkGroup.DELETE(":id", r.authMiddleware, r.controller.DeleteBookmark)
}
