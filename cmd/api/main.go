package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bzzim/glame/config"
	"github.com/bzzim/glame/controllers"
	settingController "github.com/bzzim/glame/controllers/setting"
	"github.com/bzzim/glame/middleware"
	"github.com/bzzim/glame/models"
	"github.com/bzzim/glame/pkg/token"
	"github.com/bzzim/glame/pkg/utils"
	"github.com/bzzim/glame/pkg/weather"
	"github.com/bzzim/glame/routes"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "config/api.yml", "path to config file")
}

func main() {
	flag.Parse()
	apiConfig, err := config.NewApiConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// TODO: for prod
	//handler := slog.NewJSONHandler(os.Stdout, nil)
	//handler := slog.NewTextHandler(os.Stdout, nil)
	//logger := slog.New(handler)
	logger := slog.Default()

	initFiles()
	initConfig()

	db, err := gorm.Open(sqlite.Open("data/db.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// err = db.AutoMigrate(&models.Category{}, &models.Bookmark{}, &models.Weather{})
	// if err != nil {
	// 	panic("failed automigrate")
	// }

	jwt := token.NewJWT(apiConfig.App.Name, apiConfig.Auth.Secret)

	r := gin.Default()
	// TODO: for prod
	//r.Use(sloggin.New(logger))
	//r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	authMiddleware := middleware.Auth(jwt) // миллд проверяет авторизацию и сразу выдает 401 если она не верная
	// нужен для методов, которые могут быть анонимные и авторизованные
	checkAuthMiddleware := middleware.Check(jwt) // миддл проверят авторизацию, но даже при неудачной прокидывает дальше записывая флаг в контекст

	weatherService := weather.WeatherApi{}

	router := r.Group("/api")
	authRouter := routes.NewAuthRouteController(controllers.NewAuthController(jwt, apiConfig.Auth.Password))
	authRouter.AuthRoute(router)

	settingLogger := logger.With(slog.Group("controller", slog.String("name", "settingController")))
	settingRouter := routes.NewRouteSettingController(
		settingController.NewController(db, settingLogger),
		authMiddleware,
		checkAuthMiddleware)
	settingRouter.AddRoutes(router)

	weatherLogger := logger.With(slog.Group("controller", slog.String("name", "weatherController")))
	weatherRouter := routes.NewWeatherRouteController(controllers.NewWeatherController(db, &weatherService, weatherLogger))
	weatherRouter.AddRoutes(router)

	appLogger := logger.With(slog.Group("controller", slog.String("name", "appController")))
	appRouter := routes.NewRouteAppController(
		controllers.NewAppController(db, appLogger),
		authMiddleware,
		checkAuthMiddleware)
	appRouter.AddRoutes(router)

	categoryLogger := logger.With(slog.Group("controller", slog.String("name", "categoryController")))
	categoryRouter := routes.NewRouteCategoryController(
		controllers.NewCategoryController(db, categoryLogger),
		authMiddleware,
		checkAuthMiddleware)
	categoryRouter.AddRoutes(router)

	bookmarkLogger := logger.With(slog.Group("controller", slog.String("name", "bookmarkController")))
	bookmarkRouter := routes.NewRouteBookmarkController(
		controllers.NewBookmarkController(db, bookmarkLogger),
		authMiddleware,
		checkAuthMiddleware)
	bookmarkRouter.AddRoutes(router)

	staticRouter := routes.NewRouteStaticController()
	staticRouter.AddRoutes(r)

	server := &http.Server{
		Addr:           ":7777",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func initFiles() {
	initialFile := "initialData/initialFiles.json"
	content, err := os.ReadFile(initialFile)
	if err != nil {
		return
	}

	var payload models.Files
	err = json.Unmarshal(content, &payload)
	if err != nil {
		fmt.Println(err)
	}
	for _, file := range payload.Files {
		var content string
		if file.IsJSON {
			j, err := json.Marshal(file.Template)
			if err == nil {
				content = string(j)
			}
		} else {
			content = file.Template.(string)
		}
		err = utils.WriteToFileIfNotExists(content, "data/"+file.Name)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func initConfig() {
	configFile := "data/config.json"
	configInitialFile := "initialData/initialConfig.json"

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		if err = utils.CopyFile(configInitialFile, configFile); err != nil {
			fmt.Println(err)
		}
	}

	configContent, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	configInitialContent, err := os.ReadFile(configInitialFile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var configJson map[string]interface{}
	err = json.Unmarshal(configContent, &configJson)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var initialJson map[string]interface{}
	err = json.Unmarshal(configInitialContent, &initialJson)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	for k, v := range initialJson {
		if _, ok := configJson[k]; !ok {
			configJson[k] = v
		}
	}

	content, err := json.MarshalIndent(configJson, "", "  ")
	if err != nil {
		log.Fatal("Marshal: ", err)
	}

	err = utils.WriteToFileIfNotExists(string(content), configFile)
	if err != nil {
		fmt.Println(err)
	}
}
