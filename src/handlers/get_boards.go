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

	GetBoardsUseCase interface {
		Handle(ctx context.Context, cmd getboards.Query) ([]domain.Board, error)
	}
)

// @Summary Get boards by User-id
// @Schemes
// @Tags boards
// @Accept json
// @Produce json
// @Param User-id path string true "User-id in uuid-format"
// @Success 200 {object}  GetBoardsResponse
// @Failure     400,404,408,500,503  {object}  ErrorResponse
// @Router /v1/boards [GET]

func (h *HttpHandler) GetBoards(c *gin.Context) {
	const op = "handlers.GetBoards"
	log := slog.Default()
	log.With("op", op)

	//TODO принимать offset и limtit из query параметров
	userID := c.GetHeader("User-ID")

	cmd, err := getboards.NewQuery(userID)
	if err != nil {
		log.Warn("failed to create query",
			slog.String("err", err.Error()),
			slog.String("User-id", userID),
		)
		NewErrorResponse(c, http.StatusBadRequest, "failed to create command")
		return
	}

	boards, err := h.getBoardsUC.Handle(c.Request.Context(), cmd)
	if err != nil {
		log.Error("failed get boards",
			slog.String("err", err.Error()),
			slog.String("User-id", userID),
		)
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
