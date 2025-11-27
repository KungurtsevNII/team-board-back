package main

import (
	_ "embed"
	"fmt"
	_ "io/fs"
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"github.com/KungurtsevNII/team-board-back/src/config"
	"github.com/KungurtsevNII/team-board-back/src/handlers"
	"github.com/KungurtsevNII/team-board-back/src/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/KungurtsevNII/team-board-back/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	mainPath = "/api"
)

type HttpServer struct {
	router   *gin.Engine
	srv      *http.Server
	handlers *handlers.HttpHandler
	cfg      *config.Config
}

// @title           Team Board API
// @version         1.0
// @description     API для проекта Team Board
// @host      localhost:8080
// @BasePath  /api
// @securityDefinitions.basic  BasicAuth
// @externalDocs.description  OpenAPI
func initAndStartHTTPServer(
	cfg *config.Config,
	handlers *handlers.HttpHandler,
) (*HttpServer, <-chan error) {
	log := slog.Default()
	const op = "initAndStartHttpServer"

	httpErrCh := make(chan error)

	router := gin.Default()
	router.Use(middlewares.RequestLogger()) //Логирование запросов до основной ручки

	//TODO: Поменять AllowOrigins: []string{"*"}, на хост фронта
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                // Разрешенные источники
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},            // Разрешенные методы
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "token"}, // Разрешенные заголовки
		ExposeHeaders:    []string{"Content-Length"},                                   // Заголовки, которые могут быть доступны клиенту
		AllowCredentials: true,                                                         // Разрешить отправку учетных данных (например, куки)
		MaxAge:           12 * time.Hour,                                               // Время кэширования preflight-запросов
	}))

	docs.SwaggerInfo.BasePath = mainPath
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HttpConfig.Port),
		Handler: router,
	}

	s := &HttpServer{
		cfg:      cfg,
		srv:      srv,
		handlers: handlers,
		router:   router,
	}

	//TODO: Добавить рекавери

	mainGroup := router.Group(mainPath)

	mainGroup.GET("/healthcheck", handlers.Healthcheck)

	v1Group := mainGroup.Group("/v1")
	{
		v1Group.POST("/boards/:board_id/columns", handlers.CreateColumn)
		v1Group.DELETE("/columns/:column_id", handlers.DeleteColumn)
		v1Group.POST("/boards", handlers.CreateBoard)
		v1Group.POST("/tasks", handlers.CreateTask)
		v1Group.GET("/tasks/:task_id", handlers.GetTask)
		v1Group.DELETE("/tasks/:task_id", handlers.DeleteTask)
		v1Group.GET("/boards", handlers.GetBoards)
		v1Group.GET("/boards/:id", handlers.GetBoard)
	}

	log.Info("http server is running", slog.String("port", strconv.Itoa(cfg.HttpConfig.Port)),
		slog.String("port", strconv.Itoa(cfg.HttpConfig.Port)),
		slog.String("swagger", fmt.Sprintf("http://localhost:%d/swagger/index.html", cfg.HttpConfig.Port)))

	go func() {
		if err := s.router.Run(":" + strconv.Itoa(cfg.HttpConfig.Port)); err != nil {
			httpErrCh <- errors.Wrap(err, op)
		}
	}()

	return s, httpErrCh
}
