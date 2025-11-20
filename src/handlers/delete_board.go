package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/usecase/deleteboard"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	DeleteBoardRequest struct {
		BoardID string `json:"board_id"`
	}

	DeleteBoardResponse struct {
		Success bool `json:"success"`
	}

	DeleteBoardUseCase interface {
		Handle(ctx context.Context, cmd deleteboard.Command) error
	}
)

func (h *HttpHandler) DeleteBoard(c *gin.Context) {
	const op = "handlers.DeleteBoards"
	log := slog.Default()
	log.With("op", op)

	var req DeleteBoardRequest
	if err := c.BindJSON(&req); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return
	}

	uid, err := uuid.Parse(req.BoardID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid board id")
		return
	}
	cmd, err := deleteboard.NewCommand(uid)
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
		}
		return
	}

	c.JSON(http.StatusOK, DeleteBoardResponse{Success: true})
}
