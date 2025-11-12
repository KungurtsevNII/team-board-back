package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	_ "io/fs"
	"net/http"
	"strconv"
	"time"

	"log/slog"

	"github.com/KungurtsevNII/team-board-back/src/config"
	"github.com/KungurtsevNII/team-board-back/src/handlers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
)

//go:embed openapi.yaml
var openAPISpec []byte

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
	router.Use(RequestLogger()) //Логирование запросов до основной ручки

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

	//Загрузка swagger
	router.GET("/openapi.yaml", func(c *gin.Context) {
        c.Data(http.StatusOK, "application/yaml", openAPISpec)
    })
    router.GET("/docs/*any", ginSwagger.WrapHandler(
        swaggerFiles.Handler,
        ginSwagger.URL("/openapi.yaml"),
    ))

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
		v1Group.POST("/board", handlers.CreateBoard)
	}

	log.Info("http server is running", slog.String("port", strconv.Itoa(cfg.HttpConfig.Port)),
		slog.String("port", strconv.Itoa(cfg.HttpConfig.Port)),
		slog.String("swagger", fmt.Sprintf("http://localhost:%d/docs/index.html", cfg.HttpConfig.Port)))

	go func() {
		if err := s.router.Run(":" + strconv.Itoa(cfg.HttpConfig.Port)); err != nil {
			httpErrCh <- fmt.Errorf("%s: %w", op, err)
		}
	}()

	return s, httpErrCh
}


func RequestLogger() gin.HandlerFunc {
    return func(c *gin.Context) {
		log := slog.Default()
        body, _ := io.ReadAll(c.Request.Body)
        c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

        headers := c.Request.Header

        query := c.Request.URL.Query()

        pathParams := make(map[string]string)
        for _, p := range c.Params {
            pathParams[p.Key] = p.Value
        }

        log.Info(fmt.Sprintf("%s %s",c.Request.Method, c.Request.URL.Path),
            "headers", headers,
            "query", query,
            "params", pathParams,
            "path", c.Request.URL.Path,
            "body", string(body),
        )

        c.Next()
    }
}