package handlers

import (
	"log/slog"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/config"
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	cfg            *config.HTTPConfig
	createColumnUC CreateColumnUseCase
	createBoardUC  CreateBoardUseCase
	getBoardUC     GetBoardUseCase
	createTaskUC   CreateTaskUseCase
	getBoardsUC    GetBoardsUseCase
	deleteboardUC  DeleteBoardUseCase
	getTaskUC      GetTaskUseCase
	deleteTaskUC   DeleteTaskUseCase
	deleteColumnUC DeleteColumnUseCase
	searchTasksUC      SearchTasksUseCase
}

func NewHttpHandler(
	cfg *config.HTTPConfig,
	createColumnUC CreateColumnUseCase,
	createBoardUC CreateBoardUseCase,
	getBoardUC GetBoardUseCase,
	createTaskUC CreateTaskUseCase,
	getboardsUC GetBoardsUseCase,
	deleteboardUC DeleteBoardUseCase,

	getTaskUC GetTaskUseCase,
	deleteTaskUC DeleteTaskUseCase,
	deleteColumnUC DeleteColumnUseCase,
	searchTasksUC SearchTasksUseCase,
) *HttpHandler {
	return &HttpHandler{
		cfg:            cfg,
		createColumnUC: createColumnUC,
		createBoardUC:  createBoardUC,
		getBoardUC:     getBoardUC,
		createTaskUC:   createTaskUC,
		getBoardsUC:    getboardsUC,
		deleteboardUC:  deleteboardUC,
		getTaskUC:      getTaskUC,
		deleteTaskUC:   deleteTaskUC,
		deleteColumnUC: deleteColumnUC,
		searchTasksUC:  searchTasksUC,
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
