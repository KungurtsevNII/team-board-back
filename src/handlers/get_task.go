package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/gettask"
	"github.com/gin-gonic/gin"
)

type (
	GetTaskResponse struct {
		ID          string         `json:"id"`
		ColumnID    string         `json:"column_id"`
		BoardID     string         `json:"board_id"`
		Number      int64          `json:"number"`
		Title       string         `json:"title"`
		Description *string        `json:"description"`
		Tags        []string       `json:"tags"`
		Checklists  []ChecklistDto `json:"checklists"`
		CreatedAt   time.Time      `json:"created_at"`
		UpdatedAt   time.Time      `json:"updated_at"`
		DeletedAt   *time.Time     `json:"deleted_at"`
	}

	GetTaskUseCase interface {
		Handle(ctx context.Context, query gettask.GetTaskQuery) (*domain.Task, error)
	}
)
// @Summary Получение задачи по ID
// @Schemes
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task_id path string true "ID задачи"
// @Success 200 {object}  GetTaskResponse
// @Failure     400,404,408,500,503  {object}  ErrorResponse
// @Router /v1/tasks/{task_id} [GET]
func (h *HttpHandler) GetTask(c *gin.Context) {
	const op = "handlers.GetTask"
	log := slog.Default()
	log.With("op", op)

	taskID := c.Param("task_id")

	cmd, err := gettask.NewGetTaskQuery(taskID)
	if err != nil {
		log.Warn("failed to create command", "error", err)
		switch {
		case errors.Is(err, gettask.ErrInvalidTaskID):
			NewErrorResponse(c, http.StatusBadRequest, "invalid task id")
		default:
			NewErrorResponse(c, http.StatusBadRequest, "failed to create command")
		}
		return
	}

	task, err := h.getTaskUC.Handle(c.Request.Context(), cmd)
	if err != nil {
		log.Error("failed to handle task", "error", err)
		switch {
		case errors.Is(err, gettask.ErrTaskNotFound):
			NewErrorResponse(c, http.StatusNotFound, "task not found")
		case errors.Is(err, gettask.ErrGetTaskUnknown):
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

	resp := taskDomainToGetTaskResponse(task)

	c.JSON(http.StatusOK, resp)
}

func taskDomainToGetTaskResponse(task *domain.Task) *GetTaskResponse {
	checklistResp := make([]ChecklistDto, 0, len(task.Checklists))
	for _, checklist := range task.Checklists {
		checklistItemsResp := make([]CheckListItemDto, 0, len(checklist.Items))
		for _, item := range checklist.Items {
			checklistItemsResp = append(checklistItemsResp, CheckListItemDto{
				Title:     item.Title,
				Completed: item.Completed,
			})
		}
		checklistResp = append(checklistResp, ChecklistDto{
			Title: checklist.Title,
			Items: checklistItemsResp,
		})
	}
	return &GetTaskResponse{
		ID:          task.ID.String(),
		ColumnID:    task.ColumnID.String(),
		BoardID:     task.BoardID.String(),
		Number:      task.Number,
		Title:       task.Title,
		Description: task.Description,
		Tags:        task.Tags,
		Checklists:  checklistResp,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		DeletedAt:   task.DeletedAt,
	}
}