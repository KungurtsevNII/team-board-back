package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"github.com/KungurtsevNII/team-board-back/src/config"
	"github.com/KungurtsevNII/team-board-back/src/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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

func initAndStartHTTPServer(
	cfg *config.Config,
	handlers *handlers.HttpHandler,
) (*HttpServer, <-chan error) {
	log := slog.Default()
	const op = "initAndStartHttpServer"

	httpErrCh := make(chan error)

	router := gin.Default()
	//Это нужно для того чтобы фронт мог достучаться, пока AllowOrigins: []string{"*"}, но потом это нужно поменять на хост фронта
	//TODO: Поменять AllowOrigins: []string{"*"}, на хост фронта
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                                // Разрешенные источники
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},            // Разрешенные методы
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "token"}, // Разрешенные заголовки
		ExposeHeaders:    []string{"Content-Length"},                                   // Заголовки, которые могут быть доступны клиенту
		AllowCredentials: true,                                                         // Разрешить отправку учетных данных (например, куки)
		MaxAge:           12 * time.Hour,                                               // Время кэширования preflight-запросов
	}))

	//TODO: Добавить сваггер
	// Будет тут будет запуск сваггера
	// docs.SwaggerInfo.BasePath = mainPath
	// routerGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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
		v1Group.POST("/columns", handlers.CreateColumn)
		v1Group.GET("/columns/:id", handlers.GetColumn)
		v1Group.PUT("/board", handlers.CreateBoard)
	}

	log.Info("http server is running", slog.String("port", strconv.Itoa(cfg.HttpConfig.Port)),
		slog.String("port", strconv.Itoa(cfg.HttpConfig.Port)),
		slog.String("swagger", fmt.Sprintf("http://localhost:%d/swagger/index.html", cfg.HttpConfig.Port)))

	go func() {
		if err := s.router.Run(":" + strconv.Itoa(cfg.HttpConfig.Port)); err != nil {
			httpErrCh <- fmt.Errorf("%s: %w", op, err)
		}
	}()

	return s, httpErrCh
}
