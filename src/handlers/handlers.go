package handlers

import (
	"log/slog"
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
	createBoardUC  CreateBoardUseCase
	createTaskUC   CreateTaskUseCase
	getBoardsUC    GetBoardsUseCase
	getTaskUC      GetTaskUseCase
	moveTaskUC     MoveTaskUseCase
}

func NewHttpHandler(
	cfg *config.HTTPConfig,
	createColumnUC CreateColumnUseCase,
	createBoardUC CreateBoardUseCase,
	createTaskUC CreateTaskUseCase,
	getboardsUC GetBoardsUseCase,
	getTaskUC GetTaskUseCase,
	moveTaskUC MoveTaskUseCase,
) *HttpHandler {
	return &HttpHandler{
		cfg:            cfg,
		createColumnUC: createColumnUC,
		createBoardUC:  createBoardUC,
		createTaskUC:   createTaskUC,
		getBoardsUC:    getboardsUC,
		getTaskUC:      getTaskUC,
		moveTaskUC:     moveTaskUC,
	}
}

type ErrorResponse struct {
	Err Error `json:"error"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	err := ErrorResponse{
		Err: Error{
			Code:    statusCode,
			Message: message,
		},
	}
	c.AbortWithStatusJSON(statusCode, err)
}

func (s *HttpHandler) Healthcheck(c *gin.Context) {
	const op = "handlers.Healthcheck"
	log := slog.Default().With("op", op)
	log.Info("healthcheck endpoint called")

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}
