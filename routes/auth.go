package routes

import (
	"github.com/bzzim/glame/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	controller controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (r *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("auth")
	router.POST("", r.controller.Login)
	router.POST("/validate", r.controller.Validate)
}
