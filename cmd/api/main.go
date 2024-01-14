package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bzzim/glame/config"
	"github.com/bzzim/glame/helper"
	"github.com/bzzim/glame/middleware"
	"github.com/bzzim/glame/pkg/token"
	"github.com/bzzim/glame/routes"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
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
		log.Fatal("Can't read config: ", err)
	}

	var handler slog.Handler
	if apiConfig.App.Debug {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelWarn,
		})
		gin.SetMode(gin.ReleaseMode)
	}

	logger := slog.New(handler)

	helper.InitFiles()
	helper.InitConfig()

	db, err := gorm.Open(sqlite.Open("data/db.sqlite"), &gorm.Config{})
	if err != nil {
		log.Fatal("Can't connect to database: ", err)
	}

	r := gin.New()
	if apiConfig.App.Debug {
		r.Use(gin.Logger())
	} else {
		r.Use(sloggin.New(logger))
	}
	jwt := token.NewJWT(apiConfig.App.Name, apiConfig.Auth.Secret)
	r.Use(gin.Recovery(), middleware.Cors(), middleware.CheckAuth(jwt, logger))

	routes.NewApiRoute(r, apiConfig, jwt, db, logger)

	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", apiConfig.App.Port),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	logger.Info("Start sever", slog.String("addr", server.Addr))

	if err = server.ListenAndServe(); err != nil {
		log.Fatal("Can't start server", err)
	}
}
