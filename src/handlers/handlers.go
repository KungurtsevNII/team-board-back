package handlers

import (
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/config"
	"github.com/KungurtsevNII/team-board-back/src/repository"
	"github.com/KungurtsevNII/team-board-back/src/usecase/createcolumn"
	"github.com/gin-gonic/gin"
)

const (
	v1 = "/v1"
)

type HttpHandler struct {
	cfg      *config.HTTPConfig
	columnUC ColumnUseCaseInf
	// taskUC   TaskUseCase
	// boardUC  BoardUseCase
}

func NewHttpHandler(
	cfg *config.HTTPConfig,
	repo repository.RepositoryInf, //будет прокидываться в драйверы tasks , boards , columns
) *HttpHandler {

	return &HttpHandler{
		cfg: cfg,

		columnUC: &ColumnUseCase{
			repo: repo,
			createcolumn: &createcolumn.UC{
				Repo: repo,
			},
		},

		// taskUC:         taskUC,
		// boardUC:        boardUC,
	}
}

// type BoardUseCase interface {

// }

// type TaskUseCase interface {

// }

// сюад добавляем каждый usecase , эта структура выступает как прослойка между интерфейсом бд и usecase`ами
type ColumnUseCase struct {
	repo repository.RepositoryInf

	createcolumn *createcolumn.UC
}

type ColumnUseCaseInf interface {
	CreateColumnHandle(cmd createcolumn.CreateColumnCommand) error
}

func (s *HttpHandler) Healthcheck(c *gin.Context) {
	// const op = "handlers.Healthcheck"
	// log := s.log.With("op", op, "method", c.Request.Method)
	// log.Info(c.Request.URL.Path)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
