package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"github.com/KungurtsevNII/team-board-back/src/config"
	"github.com/KungurtsevNII/team-board-back/src/handlers"
	"github.com/KungurtsevNII/team-board-back/src/repository/postgres"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
    router *gin.Engine
	srv *http.Server
	rep *postgres.Repository
	cfg *config.Config
}


func initAndStartHTTPServer(
	cfg *config.Config,
	repo *postgres.Repository,
	) error {
	log := slog.Default()
	const op = "initAndStartHttpServer"

	// ch := make(chan error)
	
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
		cfg: cfg,
		srv: srv,
		router: router,
	}
	
	//Handlers

	//TODO: Добавить хендлеры
	//TODO: Добавить рекавери
	handlers := handlers.NewHttpHandler(&cfg.HttpConfig, repo)

	
	log.Info("http server is running", slog.String("port", strconv.Itoa(cfg.HttpConfig.Port)),
		slog.String("swagger", fmt.Sprintf("http://localhost:%d/swagger/index.html", cfg.HttpConfig.Port)))


	if err := s.router.Run(":" + strconv.Itoa(cfg.HttpConfig.Port)); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}