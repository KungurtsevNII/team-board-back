package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/usecase/deletecolumn"
	"github.com/gin-gonic/gin"
)

type (
	DeleteColumnUseCase interface {
		Handle(ctx context.Context, cmd deletecolumn.Command) error
	}
)

// @Summary Удаление колонки по id
// @Schemes
// @Tags Columns
// @Accept json
// @Produce json
// @Param column_id path string true "ID колонки"
// @Success 204
// @Failure     400,404,408,409,500,503  {object}  ErrorResponse
// @Router /v1/columns/{column_id} [DELETE]
func (h *HttpHandler) DeleteColumn(c *gin.Context) {
	const op = "handlers.DeleteColumn"
	log := slog.Default()
	log.With("op", op)

	colID := c.Param("column_id")

	cmd, err := deletecolumn.NewCommand(colID)
	if err != nil {
		log.Warn("failed to create command", "error", err)
		NewErrorResponse(c, http.StatusBadRequest, "failed to create command")
		return
	}

	if err := h.deleteColumnUC.Handle(c.Request.Context(), cmd); err != nil {
		log.Error("failed to handle column", "error", err)
		switch {
		case errors.Is(err, deletecolumn.ErrDeleteColumnUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to delete column")
		case errors.Is(err, deletecolumn.ErrColumnNotFound):
			NewErrorResponse(c, http.StatusNotFound, "column not found")
		case errors.Is(err, deletecolumn.ErrGetColumnUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to get column")
		case errors.Is(err, deletecolumn.ErrCheckColumnIsEmptyUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to checking if column is empty")
		case errors.Is(err, deletecolumn.ErrColumnNotEmpty):
			NewErrorResponse(c, http.StatusConflict, "column is not empty") //Не уверен что статус 409, но вроде подходит
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
