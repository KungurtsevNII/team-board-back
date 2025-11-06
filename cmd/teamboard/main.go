package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/KungurtsevNII/team-board-back/src/config"
	"github.com/KungurtsevNII/team-board-back/src/repository/postgres"
	"github.com/sytallax/prettylog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"

)

var(
	httpError chan error
)

func main() {
	cfg := config.MustLoad()    //Сделал другой инит конфига
	log := setupLogger(cfg.Env) //И логгер читаемый

	log.Info("starting application", slog.String("env", cfg.Env))

	rep, err := postgres.New(cfg.StoragePath)
	if err != nil {
		// log.Error("can't create repository", slog.Error(err))
		panic(err)
	}

	httpsrv, httpErrCh := initAndStartHTTPServer(cfg, rep)
	if err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

    select {
    case err := <-httpErrCh:
        log.Error("http server failed", slog.Any("error", err))
        os.Exit(1)
    case sig := <-stop:
        log.Info("received shutdown signal", slog.String("signal", sig.String()))
		httpsrv.srv.Close()
		httpsrv.rep.Close()
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
