package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/getboards"
	"github.com/gin-gonic/gin"
)

type (
	GetBoardsResponse struct {
		Boards []domain.Board `json:"boards"`
	}

	GetBoardsRequest struct {
		UserID string `json:"user_id"`
	}

	GetBoardsUseCase interface {
		Handle(cmd getboards.GetBoardsCommand, ctx context.Context) (*[]domain.Board, error)
	}
)

func (h *HttpHandler) GetBoards(c *gin.Context) {
	const op = "handlers.GetBoards"
	log := slog.Default()
	log.With("op", op)

	//TODO получать user_id из тела запроса
	userID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

	cmd, err := getboards.NewGetBoardsCommand(userID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "failed to create command")
		return
	}

	boards, err := h.getBoardsUC.Handle(cmd, c.Request.Context())
	if err != nil {
		switch {
		case errors.Is(err, getboards.ErrInvalidUserID):
			NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	c.JSON(http.StatusOK, GetBoardsResponse{Boards: *boards})
}
