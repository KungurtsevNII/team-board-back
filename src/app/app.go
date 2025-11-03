package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/repository/postgres"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
	log  *slog.Logger
	rep *postgres.Repository
	port int
	srv  *http.Server //Нужно для graceful shutdown
}

func New(log *slog.Logger, e *gin.Engine, port int, rep *postgres.Repository) *App {
	// psqlStorage, err := rep.New(storagePath)
	// if err != nil{
	// 	log.Error("failed to create sql storage", slog.String("error", err.Error()))
	// 	return nil
	// }
	// restServer := gin.Default()
	// service := subAggService.New(log, psqlStorage, psqlStorage, psqlStorage, psqlStorage)

	// docs.SwaggerInfo.BasePath = mainPath
	// path := restServer.Group(mainPath)

	// subAggRest.Register(path, service)

	// restServer.GET("/healthcheck", func(g *gin.Context) {
	// 	g.JSON(http.StatusOK, http.NoBody)
	// })
	// restServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return &App{
		log: log,
		engine: e,
		rep: rep,
		port: port,
	}
}

//Для паники при ошибке
func (a *App) MustRun(){
    if err := a.Run(); err != nil{
		panic(err)
	}
}

func (a *App) Run() error{
	const op = "app.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	//Это нужно для того чтобы фронт мог достучаться, пока AllowOrigins: []string{"*"}, но потом это нужно поменять на хост фронта
	//TODO: Поменять AllowOrigins: []string{"*"}, на хост фронта
	a.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // Разрешенные источники
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},   // Разрешенные методы
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "token"}, // Разрешенные заголовки
		ExposeHeaders:    []string{"Content-Length"},                          // Заголовки, которые могут быть доступны клиенту
		AllowCredentials: true,                                                // Разрешить отправку учетных данных (например, куки)
		MaxAge:           12 * time.Hour,                                      // Время кэширования preflight-запросов
	}))

	a.srv = &http.Server{
		Addr:    fmt.Sprintf(":%d",a.port),
		Handler: a.engine,
	}

	log.Info("REST server is running", slog.String("port", strconv.Itoa(a.port)), 
		slog.String("swagger", fmt.Sprintf("http://localhost:%d/swagger/index.html",a.port)))

	if err := a.engine.Run(":" + strconv.Itoa(a.port)); err != nil{
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil

}

func (a *App) Stop(ctx context.Context) error{
    const op = "restapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping REST server", slog.Int("port", a.port))

	if err := a.srv.Shutdown(ctx); err != nil{
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
