package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/KungurtsevNII/team-board-back/src/app"
	"github.com/KungurtsevNII/team-board-back/src/config"
	"github.com/KungurtsevNII/team-board-back/src/handlers"
	"github.com/sytallax/prettylog"

	"github.com/gin-gonic/gin"
)

const (
	mainPath = "/api"
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()    //Сделал другой инит конфига
	log := setupLogger(cfg.Env) //И логгер читаемый

	log.Info("starting application", slog.String("env", cfg.Env))

	r := gin.Default()
	routerGroup := r.Group(mainPath)
	handlers.RegisterHandlers(
		log,
		routerGroup,
		nil, //Тут будет интерфейса TeamBoardAggregation,
		// который будет реализовывать repository (репозиторий будет разбивать большой
		//  интерфейс на подинтерфейсы и так они будут друг друга инплементить)
	)

	//TODO: Добавить сваггер
	// Будет тут будет запуск сваггера
	// docs.SwaggerInfo.BasePath = mainPath
	// routerGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	app := app.New(
		log,
		r,
		cfg.REST.Port,
		nil, //Тут будет инстанс базы данных
	)

	go app.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop
	app.Stop(context.Background())
	log.Info("stop gratefully")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		prettyHandler := prettylog.NewHandler(&slog.HandlerOptions{
			Level:       slog.LevelDebug,
			AddSource:   false,
			ReplaceAttr: nil,
		})
		log = slog.New(prettyHandler)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	slog.SetDefault(log)

	return log
}
