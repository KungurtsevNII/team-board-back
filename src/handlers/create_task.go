package handlers

import (
	"context"
	"encoding/json"
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
		ColumnID    string                   `json:"column_id"`
		BoardID     string                   `json:"board_id"`
		Title       string                   `json:"title"`
		Description *string                  `json:"description"`
		Tags        []string                 `json:"tags"`
		Checklists  []ChecklistDto `json:"checklists"`
	}

	CreateTaskResponse struct {
		ID          string                    `json:"id"`
		ColumnID    string                    `json:"column_id"`
		BoardID     string                    `json:"board_id"`
		Number      int64                     `json:"number"`
		Title       string                    `json:"title"`
		Description *string                   `json:"description"`
		Tags        []string                  `json:"tags"`
		Checklists  []ChecklistDto `json:"checklists"`
		CreatedAt   time.Time                 `json:"created_at"`
		UpdatedAt   time.Time                 `json:"updated_at"`
		DeletedAt   *time.Time                `json:"deleted_at"`
	}

	ChecklistDto struct {
		Title string `json:"title"`
		Items []struct {
			Title     string `json:"title"`
			Completed bool   `json:"completed"`
		} `json:"items"`
	}

	CreateTaskUseCase interface {
		Handle(
			ctx context.Context,
			cmd createtask.CreateTaskCommand,
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

	checklistsJSON, err := json.Marshal(req.Checklists)
	if err != nil {
		log.Warn("failed to bind request", slog.String("err", err.Error()))
		NewErrorResponse(c, http.StatusBadRequest, "bad body")
		return
	}
	rawMsg := json.RawMessage(checklistsJSON)

	cmd, err := createtask.NewCreateTaskCommand(
		req.ColumnID,
		req.BoardID,
		req.Title,
		req.Description,
		req.Tags,
		&rawMsg,
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


	resp := CreateTaskResponse{
		ID:          dmn.ID.String(),
		ColumnID:    dmn.ColumnID.String(),
		BoardID:     dmn.BoardID.String(),
		Number:      dmn.Number,
		Title:       dmn.Title,
		Description: dmn.Description,
		Tags:        dmn.Tags,
		Checklists:  req.Checklists,
		CreatedAt:   dmn.CreatedAt,
		UpdatedAt:   dmn.UpdatedAt,
		DeletedAt:   dmn.DeletedAt,
	}

	c.JSON(http.StatusCreated, resp)
}
