package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/getboards"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	Board struct {
		ID        uuid.UUID `db:"id" json:"id"`
		Name      string    `db:"name" json:"name"`
		ShortName string    `db:"short_name" json:"short_name"`
		UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	}

	GetBoardsResponse struct {
		Boards []Board `json:"boards"`
	}

	GetBoardsRequest struct {
		UserID string `json:"user_id"`
	}

	GetBoardsUseCase interface {
		Handle(ctx context.Context, cmd getboards.Query) ([]domain.Board, error)
	}
)

func (h *HttpHandler) GetBoards(c *gin.Context) {
	const op = "handlers.GetBoards"
	log := slog.Default()
	log.With("op", op)

	//TODO принимать offset и limtit из query параметров
	//TODO получать user_id из тела запроса
	userID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"

	cmd, err := getboards.NewQuery(userID)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "failed to create command")
		return
	}

	boards, err := h.getBoardsUC.Handle(c.Request.Context(), cmd)
	if err != nil {
		switch {
		case errors.Is(err, getboards.ErrInvalidUserID):
			NewErrorResponse(c, http.StatusBadRequest, "invalid user id")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	boardsResp := make([]Board, len(boards))
	for i, board := range boards {
		boardsResp[i] = Board{
			ID:        board.ID,
			Name:      board.Name,
			ShortName: board.ShortName,
			UpdatedAt: board.UpdatedAt,
		}
	}

	c.JSON(http.StatusOK, GetBoardsResponse{Boards: boardsResp})
}
