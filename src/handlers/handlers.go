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
	createBoardUC  CreateBoardUseCase
}

func NewHttpHandler(
	cfg *config.HTTPConfig,
	columnUC CreateColumnUseCase,
	getColumnUC GetColumnUseCase,
	createBoardUC CreateBoardUseCase,

) *HttpHandler {
	return &HttpHandler{
		cfg:            cfg,
		createColumnUC: columnUC,
		getColumnUC:    getColumnUC,
		createBoardUC:  createBoardUC,
	}
}

type ErrorResponse struct {
	Err Error `json:"error"`
}

type Error struct{
	Code int `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context,statusCode int, message string) {
	err := ErrorResponse{
		Err : Error{
			Code: statusCode,
			Message: message,
		},
	}
	c.AbortWithStatusJSON(statusCode, err)
}

func (s *HttpHandler) Healthcheck(c *gin.Context) {
	// const op = "handlers.Healthcheck"
	// log := s.log.With("op", op, "method", c.Request.Method)
	// log.Info(c.Request.URL.Path)

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
