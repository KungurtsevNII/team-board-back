package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/createtask"
)

type (
	CreateTaskRequest struct {
		ColumnID    string         `json:"column_id"`
		BoardID     string         `json:"board_id"`
		Title       string         `json:"title"`
		Description *string        `json:"description"`
		Tags        []string       `json:"tags"`
		Checklists  []ChecklistDto `json:"checklists"`
	}

	CreateTaskResponse struct {
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

	ChecklistDto struct {
		Title string             `json:"title"`
		Items []CheckListItemDto `json:"items"`
	}
	CheckListItemDto struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	CreateTaskUseCase interface {
		Handle(
			ctx context.Context,
			cmd createtask.Command,
		) (task *domain.Task, err error)
	}
)

// @Summary Создание новой задачи
// @Schemes
// @Tags Tasks
// @Accept json
// @Produce json
// @Param createTaskRequest body CreateTaskRequest true "request на создание таски"
// @Success 201 {object}  CreateTaskResponse
// @Failure     400,404,408,500,503  {object}  ErrorResponse
// @Router /v1/tasks [POST]
func (h *HttpHandler) CreateTask(c *gin.Context) {
	const op = "handlers.CreateTask"
	log := slog.Default()
	log.With("op", op)

	var req CreateTaskRequest
	if err := c.BindJSON(&req); err != nil {
		log.Warn("failed to bind request", slog.String("err", err.Error()))
		NewErrorResponse(c, http.StatusBadRequest, "bad body")
		return
	}

	checkListsDmn := make([]domain.Checklist, 0, len(req.Checklists))
	for _, checklist := range req.Checklists {
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

	cmd, err := createtask.NewCommand(
		req.ColumnID,
		req.BoardID,
		req.Title,
		req.Description,
		req.Tags,
		checkListsDmn,
	)

	if err != nil {
		log.Warn("failed to create command", slog.String("err", err.Error()))

		switch {
		case errors.Is(err, createtask.ErrValidationFailed):
			NewErrorResponse(c, http.StatusBadRequest, "validation failed")
		case errors.Is(err, createtask.ErrInvalidUUID):
			NewErrorResponse(c, http.StatusBadRequest, "invalid id")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	dmn, err := h.createTaskUC.Handle(c, cmd)
	if err != nil {
		log.Error("failed to create column", slog.String("err", err.Error()))

		switch {
		case errors.Is(err, createtask.ErrColumnOrBoardIsNotExists):
			NewErrorResponse(c, http.StatusNotFound, "board or column not found")
		case errors.Is(err, createtask.ErrGetLastNumberFailed):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to get last number")
		case errors.Is(err, createtask.ErrValidationFailed):
			NewErrorResponse(c, http.StatusBadRequest, "validation failed")
		case errors.Is(err, createtask.ErrCreateTaskUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to create task")
		case errors.Is(err, context.Canceled):
			NewErrorResponse(c, http.StatusRequestTimeout, "request canceled")
		case errors.Is(err, context.DeadlineExceeded):
			NewErrorResponse(c, http.StatusServiceUnavailable, "request timeout")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	checklistResp := make([]ChecklistDto, 0, len(dmn.Checklists))
	for _, checklist := range dmn.Checklists {
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
	resp := CreateTaskResponse{
		ID:          dmn.ID.String(),
		ColumnID:    dmn.ColumnID.String(),
		BoardID:     dmn.BoardID.String(),
		Number:      dmn.Number,
		Title:       dmn.Title,
		Description: dmn.Description,
		Tags:        dmn.Tags,
		Checklists:  checklistResp,
		CreatedAt:   dmn.CreatedAt,
		UpdatedAt:   dmn.UpdatedAt,
		DeletedAt:   dmn.DeletedAt,
	}

	c.JSON(http.StatusCreated, resp)
}
