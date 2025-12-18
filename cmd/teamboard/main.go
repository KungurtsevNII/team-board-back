package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/KungurtsevNII/team-board-back/src/config"
	"github.com/KungurtsevNII/team-board-back/src/handlers"
	"github.com/KungurtsevNII/team-board-back/src/repository/postgres"
	"github.com/KungurtsevNII/team-board-back/src/usecase/createboard"
	"github.com/KungurtsevNII/team-board-back/src/usecase/createcolumn"
	"github.com/KungurtsevNII/team-board-back/src/usecase/createtask"
	"github.com/KungurtsevNII/team-board-back/src/usecase/deleteboard"
	"github.com/KungurtsevNII/team-board-back/src/usecase/deletecolumn"
	"github.com/KungurtsevNII/team-board-back/src/usecase/deletetask"
	"github.com/KungurtsevNII/team-board-back/src/usecase/getboard"
	"github.com/KungurtsevNII/team-board-back/src/usecase/getboards"
	"github.com/KungurtsevNII/team-board-back/src/usecase/gettask"
	"github.com/KungurtsevNII/team-board-back/src/usecase/searchtasks"
	"github.com/KungurtsevNII/team-board-back/src/usecase/movetask"
	"github.com/sytallax/prettylog"
	"github.com/KungurtsevNII/team-board-back/src/usecase/puttask"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.String("env", cfg.Env))
	log.Info("config", slog.Any("cfg", cfg))

	rep, err := postgres.New(cfg.PostgresConfig.Host)
	if err != nil {
		panic(err)
	}

	handlers := handlers.NewHttpHandler(
		&cfg.HttpConfig,
		createcolumn.NewUC(rep),
		createboard.NewUC(rep),
		getboard.NewUC(rep),
		createtask.NewUC(rep),
		getboards.NewUC(rep),
		deleteboard.NewUC(rep),
		gettask.NewUC(rep),
		deletetask.NewUC(rep),
		deletecolumn.NewUC(rep),
		searchtasks.NewUC(rep),
		movetask.NewUC(rep),
		puttask.NewUC(rep),
	)

	log.Info("repository connected", slog.String("path", cfg.PostgresConfig.Host))

	httpsrv, httpErrCh := initAndStartHTTPServer(cfg, handlers)
	httpsrv.state.Ready = true //Все зависимости поднялись, поэтому реди

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-httpErrCh:
		log.Error("http server failed", slog.Any("error", err))
		os.Exit(1)
	case sig := <-stop:
		log.Info("received shutdown signal", slog.String("signal", sig.String()))
		_ = httpsrv.srv.Close()
		rep.Close()
		log.Info("shutdown complete")
	}
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
