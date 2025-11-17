package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/getboard"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type (
	GetBoardRequest struct {
		ID uuid.UUID `json:"id"`
	}

	GetBoardBoard struct {
		ID        uuid.UUID        `json:"id"`
		Name      string           `json:"name"`
		ShortName string           `json:"short_name"`
		Columns   []GetBoardColumn `json:"columns"`
	}

	GetBoardColumn struct {
		ID       uuid.UUID      `json:"id"`
		BoardID  string         `json:"board_id"`
		OrderNum int            `json:"order_num"`
		Name     string         `json:"name"`
		Tasks    []GetBoardTask `json:"tasks"`
	}

	GetBoardTask struct {
		ID       uuid.UUID `json:"id"`
		ColumnID string    `json:"column_id"`
		Title    string    `json:"title"`
		Content  string    `json:"content"`
	}

	GetBoardUseCase interface {
		Handle(ctx context.Context, cmd getboard.GetBoardCommand) (domain.Board, error)
	}
)

func (h *HttpHandler) GetBoard(c *gin.Context) {
	const op = "handlers.GetBoard"
	log := slog.Default()
	log.With("op", op, "method", c.Request.Method)
	log.Info(c.Request.URL.Path)

	id := c.Param("id")
	cmd, err := getboard.NewGetBoardCommand(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	board, err := h.getBoardUC.Handle(c.Request.Context(), cmd)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err, domain.ErrBoardNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err, getboard.ErrBoardIsNotExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		}
	}

	columns := make([]GetBoardColumn, len(board.Columns))
	for i, col := range board.Columns {
		columns[i] = GetBoardColumn{
			ID:       col.ID,
			BoardID:  col.BoardID,
			OrderNum: col.OrderNum,
			Name:     col.Name,
			Tasks:    []GetBoardTask{}, // TODO : сделать таски
		}
	}

	resp := GetBoardBoard{
		Name:      board.Name,
		ShortName: board.ShortName,
		Columns:   columns,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}
