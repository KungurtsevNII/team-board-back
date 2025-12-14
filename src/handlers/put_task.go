package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/gin-gonic/gin"
	"github.com/KungurtsevNII/team-board-back/src/usecase/puttask"
)

type (
	PutTaskResponse struct {
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
	}

	PutTaskRequest struct {
		ColumnID    string         `json:"column_id"`
		BoardID     string         `json:"board_id"`
		Number      int64          `json:"number"`
		Title       string         `json:"title"`
		Description *string        `json:"description"`
		Tags        []string       `json:"tags"`
		Checklists  []ChecklistDto `json:"checklists"`
	}

	PutTaskUseCase interface {
		Handle(ctx context.Context, cmd puttask.Command) (*domain.Task, error)
	}
)
// @Summary Изменение задачи
// @Schemes
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task_id path string true "ID задачи"
// @Param putTaskRequest body PutTaskRequest true "put task request"
// @Success 200 {object}  PutTaskResponse
// @Failure     400,404,408,500,503  {object}  ErrorResponse
// @Router /v1/tasks/{task_id} [PUT]
func (h *HttpHandler) PutTask(c *gin.Context) {
	const op = "handlers.GetTask"
	log := slog.Default()
	log.With("op", op)

	taskID := c.Param("task_id")

	var req PutTaskRequest
	if err := c.BindJSON(&req); err != nil {
		log.Warn("failed to bind request", slog.String("err", err.Error()))
		NewErrorResponse(c, http.StatusBadRequest, "bad body")
		return
	}

	cmd, err := puttask.NewCommand(
		taskID,
		req.BoardID, 
		req.ColumnID, 
		req.Title, 
		req.Number, 
		req.Description, 
		req.Tags, 
		putChecklistsRequestToDomain(req.Checklists),
	)
	if err != nil {
		log.Warn("failed to create command", "error", err)
		NewErrorResponse(c, http.StatusBadRequest, "failed to create command")
		return
	}

	task, err := h.putTaskUC.Handle(c.Request.Context(), cmd)
	if err != nil {
		log.Error("failed to handle task", "error", err)
		switch {
		case errors.Is(err, puttask.ErrTaskNotFound):
			NewErrorResponse(c, http.StatusNotFound, "task not found")
		case errors.Is(err, puttask.ErrPutTaskUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to put task")
		case errors.Is(err, puttask.ErrColumnNotFound):
			NewErrorResponse(c, http.StatusNotFound, "column not found")
		case errors.Is(err, context.Canceled):
			NewErrorResponse(c, http.StatusRequestTimeout, "request canceled")
		case errors.Is(err, context.DeadlineExceeded):
			NewErrorResponse(c, http.StatusServiceUnavailable, "request timeout")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	resp := taskDomainToPutTaskResponse(task)

	c.JSON(http.StatusOK, resp)
}

func putChecklistsRequestToDomain(checklists []ChecklistDto) []domain.Checklist {
	checkListsDmn := make([]domain.Checklist, 0, len(checklists))
	for _, checklist := range checklists {
		checklistItemsDmn := make([]domain.ChecklistItem, 0, len(checklist.Items))
		for _, item := range checklist.Items {
			checklistItemsDmn = append(checklistItemsDmn, domain.NewChecklistItem(
				item.Title,
				item.Completed,
			))
		}
		checkListsDmn = append(checkListsDmn, domain.NewChecklist(
			checklist.Title,
			checklistItemsDmn,
		))
	}
	return checkListsDmn
}

func taskDomainToPutTaskResponse(task *domain.Task) *PutTaskResponse {
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
	return &PutTaskResponse{
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
	}
}