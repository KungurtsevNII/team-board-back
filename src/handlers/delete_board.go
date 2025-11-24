package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/usecase/deleteboard"
	"github.com/gin-gonic/gin"
)

type (
	DeleteBoardRequest struct {
		ID string `json:"id"`
	}

	DeleteBoardUseCase interface {
		Handle(ctx context.Context, cmd deleteboard.Command) error
	}
)

// DeleteBoard godoc
// @Summary Удаление доски
// @Description Удаляет доску по её ID
// @Tags boards
// @Accept json
// @Produce json
// @Param request body DeleteBoardRequest true "ID доски"
// @Success 204 "Доска успешно удалена"
// @Failure 400 {object} ErrorResponse "Некорректный запрос или неверный ID"
// @Failure 404 {object} ErrorResponse "Доска не найдена"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка сервера"
// @Router /v1/boards [delete]

func (h *HttpHandler) DeleteBoard(c *gin.Context) {
	const op = "handlers.DeleteBoards"
	log := slog.Default()
	log.With("op", op)

	var req DeleteBoardRequest
	if err := c.BindJSON(&req); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	cmd, err := deleteboard.NewCommand(req.ID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "failed to create command")
		return
	}

	err = h.deleteboardUC.Handle(c.Request.Context(), cmd)
	if err != nil {
		switch {
		case errors.Is(err, deleteboard.ErrBoardIdEmpty):
			NewErrorResponse(c, http.StatusBadRequest, "board id is empty")
		case errors.Is(err, deleteboard.ErrBoardIdInvalid):
			NewErrorResponse(c, http.StatusBadRequest, "board id is invalid")
		case errors.Is(err, deleteboard.ErrBoardDoesntExist):
			NewErrorResponse(c, http.StatusNotFound, "board doesn't exist")
		default:
			log.Error("failed to delete board", "error", err)
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.JSON(http.StatusNoContent, "")
}
