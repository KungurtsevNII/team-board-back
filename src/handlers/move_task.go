package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/movetask"
	"github.com/gin-gonic/gin"
)

type (
	MoveTaskRequest struct {
		ColumnID string `json:"column_id" binding:"required"`
	}

	MoveTaskResponse struct {
		ID       string `json:"id"`
		ColumnID string `json:"column_id"`
		BoardID  string `json:"board_id"`
	}

	MoveTaskUseCase interface {
		Handle(
			ctx context.Context,
			cmd movetask.MoveTaskCommand,
		) (task *domain.Task, err error)
	}
)

// @Summary Перемещение задачи в другую колонку
// @Schemes
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task_id path string true "ID задачи"
// @Param moveTaskRequest body MoveTaskRequest true "request на перемещение задачи"
// @Success 200 {object}  MoveTaskResponse
// @Failure     400,404,408,500,503  {object}  ErrorResponse
// @Router /v1/tasks/{task_id}/move [PUT]
func (h *HttpHandler) MoveTask(c *gin.Context) {
	const op = "handlers.MoveTask"
	log := slog.Default()
	log.With("op", op)

	taskID := c.Param("task_id")

	var req MoveTaskRequest
	if err := c.BindJSON(&req); err != nil {
		log.Warn("failed to bind request", slog.String("err", err.Error()))
		NewErrorResponse(c, http.StatusBadRequest, "bad body")
		return
	}

	cmd, err := movetask.NewMoveTaskCommand(taskID, req.ColumnID)
	if err != nil {
		log.Warn("failed to create command",
			slog.String("err", err.Error()),
			slog.String("task_id", taskID),
			slog.String("column_id", req.ColumnID))

		switch {
		case errors.Is(err, movetask.ErrValidationFailed):
			NewErrorResponse(c, http.StatusBadRequest, "validation failed")
		case errors.Is(err, movetask.ErrInvalidUUID):
			NewErrorResponse(c, http.StatusBadRequest, "invalid task or column id")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	dmn, err := h.moveTaskUC.Handle(c, cmd)
	if err != nil {
		log.Error("failed to move task",
			slog.String("err", err.Error()),
			slog.String("task_id", cmd.TaskID.String()),
			slog.String("column_id", cmd.ColumnID.String()))

		switch {
		case errors.Is(err, movetask.ErrTaskNotFound):
			NewErrorResponse(c, http.StatusNotFound, "task not found")
		case errors.Is(err, movetask.ErrColumnNotFound):
			NewErrorResponse(c, http.StatusNotFound, "column not found")
		case errors.Is(err, movetask.ErrColumnNotInBoard):
			NewErrorResponse(c, http.StatusBadRequest, "column does not belong to task's board")
		case errors.Is(err, movetask.ErrMoveTaskUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to move task")
		case errors.Is(err, context.Canceled):
			NewErrorResponse(c, http.StatusRequestTimeout, "request canceled")
		case errors.Is(err, context.DeadlineExceeded):
			NewErrorResponse(c, http.StatusServiceUnavailable, "request timeout")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	resp := MoveTaskResponse{
		ID:       dmn.ID.String(),
		ColumnID: dmn.ColumnID.String(),
		BoardID:  dmn.BoardID.String(),
	}

	c.JSON(http.StatusOK, resp)
}
