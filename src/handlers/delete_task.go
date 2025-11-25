package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/usecase/deletetask"
	"github.com/gin-gonic/gin"
)

type (
	DeleteTaskUseCase interface {
		Handle(ctx context.Context, cmd deletetask.Command) error
	}
)

// @Summary Удаление задачи по id
// @Schemes
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task_id path string true "ID задачи"
// @Success 204
// @Failure     400,404,408,500,503  {object}  ErrorResponse
// @Router /v1/tasks/{task_id} [DELETE]
func (h *HttpHandler) DeleteTask(c *gin.Context) {
	const op = "handlers.DeleteTask"
	log := slog.Default()
	log.With("op", op)

	taskID := c.Param("task_id")

	cmd, err := deletetask.NewCommand(taskID)
	if err != nil {
		log.Warn("failed to create command", "error", err)
		NewErrorResponse(c, http.StatusBadRequest, "failed to create command")
		return
	}

	if err := h.deleteTaskUC.Handle(c.Request.Context(), cmd); err != nil{
		log.Error("failed to handle task", "error", err)
		switch {
		case errors.Is(err, deletetask.ErrDeleteTaskUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to delete task")
		case errors.Is(err, deletetask.ErrTaskNotFound):
			NewErrorResponse(c, http.StatusNotFound, "task not found")
		case errors.Is(err, deletetask.ErrGetTaskUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to get task")
		case errors.Is(err, context.Canceled):
			NewErrorResponse(c, http.StatusRequestTimeout, "request canceled")
		case errors.Is(err, context.DeadlineExceeded):
			NewErrorResponse(c, http.StatusServiceUnavailable, "request timeout")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
