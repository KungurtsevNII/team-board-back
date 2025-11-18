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
		Tasks     []GetBoardTask   `json:"tasks"`
	}

	GetBoardColumn struct {
		ID       uuid.UUID      `json:"id"`
		BoardID  uuid.UUID      `json:"board_id"`
		OrderNum int64          `json:"order_num"`
		Name     string         `json:"name"`
		Tasks    []GetBoardTask `json:"tasks"`
	}

	GetBoardTask struct {
		ID       uuid.UUID `json:"id"`
		ColumnID uuid.UUID `json:"column_id"`
		BoardID  uuid.UUID `json:"board_id"`
		Number   int64     `json:"number"`
		Title    string    `json:"title"`
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
			Tasks:    []GetBoardTask{},
		}
	}

	tasks := make([]GetBoardTask, len(board.Tasks))
	for i, task := range board.Tasks {
		tasks[i] = GetBoardTask{
			ID:       task.ID,
			ColumnID: task.ColumnID,
			BoardID:  task.BoardID,
			Number:   task.Number,
			Title:    task.Title,
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
