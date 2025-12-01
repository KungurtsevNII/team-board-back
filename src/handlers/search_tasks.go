package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/searchtasks"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	SearchTasksRequest struct {
		Query   string `json:"query"`
		Limit   uint `json:"limit"`
		Offset  uint `json:"offset"`
		Filters struct {
			Tags []string `json:"tags"`
		} `json:"filters"`
	}

	SearchTasksUseCase interface {
		Handle(ctx context.Context, q searchtasks.Query) ([]domain.Task, error)
	}

	SearchTaskResponse struct {
		ID       uuid.UUID `json:"id"`
		ColumnID uuid.UUID `json:"column_id"`
		BoardID  uuid.UUID `json:"board_id"`
		Number   int64     `json:"number"`
		Title    string    `json:"title"`
	}
)

// @Summary Поиск задач по тегам и названию
// @Schemes
// @Tags Tasks
// @Accept json
// @Produce json
// @Param searchTasksRequest body SearchTasksRequest true "request для поиска тасок"
// @Success 200 {object}  []SearchTaskResponse
// @Failure     400,408,500,503  {object}  ErrorResponse
// @Router /v1/tasks/search [POST]
func (h *HttpHandler) SearchTasks(c *gin.Context) {
	const op = "handlers.SearchTasks"
	log := slog.Default()
	log.With("op", op)

	var req SearchTasksRequest
	err := c.BindJSON(&req)
	if err != nil {
		log.Warn("failed to bind request", "error", err)
		NewErrorResponse(c, http.StatusBadRequest, "failed to bind request")
		return
	}

	qry, err := searchtasks.NewQuery(req.Filters.Tags, req.Query, req.Limit, req.Offset)
	if err != nil {
		log.Warn("failed to create command", "error", err)
		NewErrorResponse(c, http.StatusBadRequest, "failed to create command")
		return
	}

	tasks, err := h.searchTasksUC.Handle(c.Request.Context(), qry)
	if err != nil {
		log.Error("failed to handle task", "error", err)
		switch {
		case errors.Is(err, searchtasks.ErrSearchTasks):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to search tasks")
		case errors.Is(err, context.Canceled):
			NewErrorResponse(c, http.StatusRequestTimeout, "request canceled")
		case errors.Is(err, context.DeadlineExceeded):
			NewErrorResponse(c, http.StatusServiceUnavailable, "request timeout")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	resp := taskDomainsToSearchTaskResponses(tasks)

	c.JSON(http.StatusOK, resp)
}

func taskDomainsToSearchTaskResponses(tasks []domain.Task) []SearchTaskResponse {
	resps := make([]SearchTaskResponse, 0, len(tasks))
	for _, el := range tasks {
		resps = append(resps, SearchTaskResponse{
			ID:       el.ID,
			ColumnID: el.ColumnID,
			BoardID:  el.BoardID,
			Number:   el.Number,
			Title:    el.Title,
		})
	}
	return resps
}
