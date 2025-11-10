package handlers

import (
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/config"
	"github.com/gin-gonic/gin"
)

const (
	v1 = "/v1"
)

type HttpHandler struct {
	cfg            *config.HTTPConfig
	createColumnUC CreateColumnUseCase
	getColumnUC    GetColumnUseCase
}

func NewHttpHandler(
	cfg *config.HTTPConfig,
	columnUC CreateColumnUseCase,
	getColumnUC GetColumnUseCase,
) *HttpHandler {
	return &HttpHandler{
		cfg:            cfg,
		createColumnUC: columnUC,
		getColumnUC:    getColumnUC,
	}
}

// type BoardUseCase interface {

// }

// type TaskUseCase interface {

// }

func (s *HttpHandler) Healthcheck(c *gin.Context) {
	// const op = "handlers.Healthcheck"
	// log := s.log.With("op", op, "method", c.Request.Method)
	// log.Info(c.Request.URL.Path)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
