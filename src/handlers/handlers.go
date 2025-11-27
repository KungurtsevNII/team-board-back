package handlers

import (
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
	getTaskUC      GetTaskUseCase
	deleteTaskUC   DeleteTaskUseCase
}

func NewHttpHandler(
	cfg *config.HTTPConfig,
	createColumnUC CreateColumnUseCase,
	createBoardUC CreateBoardUseCase,
	getBoardUC GetBoardUseCase,
	createTaskUC CreateTaskUseCase,
	getboardsUC GetBoardsUseCase,
	getTaskUC GetTaskUseCase,
	deleteTaskUC   DeleteTaskUseCase,
) *HttpHandler {
	return &HttpHandler{
		cfg:            cfg,
		createColumnUC: createColumnUC,
		createBoardUC:  createBoardUC,
		getBoardUC:     getBoardUC,
		createTaskUC:   createTaskUC,
		getBoardsUC:    getboardsUC,
		getTaskUC:      getTaskUC,
		deleteTaskUC:   deleteTaskUC,
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
