package routes

import (
	"log/slog"

	"github.com/bzzim/glame/config"
	"github.com/bzzim/glame/controllers"
	"github.com/bzzim/glame/controllers/setting"
	"github.com/bzzim/glame/middleware"
	"github.com/bzzim/glame/pkg/token"
	"github.com/bzzim/glame/pkg/weather"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApiRoute struct {
	router         *gin.Engine
	cfg            *config.ApiConfig
	jwt            *token.JWT
	db             *gorm.DB
	log            *slog.Logger
	middleware     gin.HandlerFunc
	weatherService weather.WeatherApi
}

func NewApiRoute(router *gin.Engine, cfg *config.ApiConfig, jwt *token.JWT, db *gorm.DB, log *slog.Logger) {
	middleware := middleware.Auth()
	weatherService := weather.New()
	apiRouter := ApiRoute{
		router:         router,
		cfg:            cfg,
		jwt:            jwt,
		db:             db,
		log:            log,
		middleware:     middleware,
		weatherService: weatherService,
	}
	apiRouter.compose()
}

func (r *ApiRoute) compose() {
	r.static()

	apiGroup := r.router.Group("/api")
	r.auth(apiGroup)
	r.settings(apiGroup)
	r.weather(apiGroup)
	r.apps(apiGroup)
	r.categories(apiGroup)
	r.bookmarks(apiGroup)
}

func (r *ApiRoute) auth(api *gin.RouterGroup) {
	controller := controllers.NewAuthController(r.jwt, r.cfg.Auth.Password)
	group := api.Group("auth")
	group.POST("", controller.Login)
	group.POST("/validate", controller.Validate)
}

func (r *ApiRoute) settings(api *gin.RouterGroup) {
	controller := setting.NewController(r.db, r.log)
	configGroup := api.Group("config")
	configGroup.GET("", controller.Config)
	configGroup.GET("0/css", controller.Css)
	configGroup.PUT("", r.middleware, controller.SaveConfig)
	configGroup.PUT("0/css", r.middleware, controller.SaveCss)

	themesGroup := api.Group("themes")
	themesGroup.GET("", controller.Themes)
	themesGroup.POST("", r.middleware, controller.AddTheme)
	themesGroup.PUT(":name", r.middleware, controller.SaveTheme)
	themesGroup.DELETE(":name", r.middleware, controller.DeleteTheme)

	queriesGroup := api.Group("queries")
	queriesGroup.GET("", controller.Queries)
	queriesGroup.POST("", r.middleware, controller.AddQuery)
	queriesGroup.PUT(":prefix", r.middleware, controller.SaveQuery)
	queriesGroup.DELETE(":prefix", r.middleware, controller.DeleteQuery)

}

func (r *ApiRoute) weather(api *gin.RouterGroup) {
	controller := controllers.NewWeatherController(r.db, &r.weatherService, r.log)
	group := api.Group("weather")
	group.GET("", controller.Weather)
	group.GET("update", controller.UpdateWeather)
}

func (r *ApiRoute) apps(api *gin.RouterGroup) {
	controller := controllers.NewAppController(r.db, r.log)
	group := api.Group("apps")
	group.GET("", controller.Apps)
	group.POST("", r.middleware, controller.AddApp)
	group.PUT(":id", r.middleware, controller.SaveApp)
	group.PUT("0/reorder", r.middleware, controller.ReorderApp)
	group.DELETE(":id", r.middleware, controller.DeleteApp)
}

func (r *ApiRoute) categories(api *gin.RouterGroup) {
	controller := controllers.NewCategoryController(r.db, r.log)
	group := api.Group("categories")
	group.GET("", controller.Categories)
	group.POST("", r.middleware, controller.AddCategory)
	group.PUT(":id", r.middleware, controller.SaveCategory)
	group.PUT("0/reorder", r.middleware, controller.ReorderCategory)
	group.DELETE(":id", r.middleware, controller.DeleteCategory)
}

func (r *ApiRoute) bookmarks(api *gin.RouterGroup) {
	controller := controllers.NewBookmarkController(r.db, r.log)
	group := api.Group("bookmarks")
	group.POST("", r.middleware, controller.AddBookmark)
	group.PUT(":id", r.middleware, controller.SaveBoomark)
	group.PUT("0/reorder", r.middleware, controller.ReorderBookmark)
	group.DELETE(":id", r.middleware, controller.DeleteBookmark)
}

func (r *ApiRoute) static() {
	r.router.Static("uploads", "./data/uploads")
	r.router.StaticFile("flame.css", "./data/flame.css")
}
