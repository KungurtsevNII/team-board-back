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

	GetBoardBoard struct {
		ID        uuid.UUID        `json:"id"`
		Name      string           `json:"name"`
		ShortName string           `json:"short_name"`
		Columns   []GetBoardColumn `json:"columns"`
		Tasks     []GetBoardTask   `json:"tasks"`
	}

	GetBoardColumn struct {
		ID       uuid.UUID `json:"id"`
		BoardID  uuid.UUID `json:"board_id"`
		OrderNum int64     `json:"order_num"`
		Name     string    `json:"name"`
	}

	GetBoardTask struct {
		ID       uuid.UUID `json:"id"`
		ColumnID uuid.UUID `json:"column_id"`
		BoardID  uuid.UUID `json:"board_id"`
		Number   int64     `json:"number"`
		Title    string    `json:"title"`
	}

	GetBoardUseCase interface {
		Handle(ctx context.Context, cmd getboard.Query) (*domain.Board, error)
	}
)

// @Summary Получение доски по id
// @Schemes
// @Tags Boards
// @Accept json
// @Produce json
// @Param id path string true "User-id in uuid-format"
// @Success 200 {object}  GetBoardsResponse
// @Failure     400,404,408,500,503  {object}  ErrorResponse
// @Router /v1/boards/{id} [GET]
func (h *HttpHandler) GetBoard(c *gin.Context) {
	const op = "handlers.GetBoard"
	log := slog.Default()
	log.With("op", op, "method", c.Request.Method)
	log.Info(c.Request.URL.Path)

	id := c.Param("id")
	cmd, err := getboard.NewQuery(id)
	if err != nil {
		log.Warn("failed to create command", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	board, err := h.getBoardUC.Handle(c.Request.Context(), cmd)
	if err != nil {
		log.Error("failed to handle board", "error", err)
		switch {
		case errors.Is(err, getboard.ErrInvalidID):
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err, getboard.ErrBoardNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
		case errors.Is(err, getboard.ErrBoardIsNotExists):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	columns := dtoColumnsToResp(board.Columns)
	tasks := dtoTasksToResp(board.Tasks)

	resp := GetBoardBoard{
		ID:        board.ID,
		Name:      board.Name,
		ShortName: board.ShortName,
		Columns:   columns,
		Tasks:     tasks,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

func dtoColumnsToResp(dtoCol []domain.Column) []GetBoardColumn {
	columns := make([]GetBoardColumn, len(dtoCol))
	for i, col := range dtoCol {
		columns[i] = GetBoardColumn{
			ID:       col.ID,
			BoardID:  col.BoardID,
			OrderNum: col.OrderNum,
			Name:     col.Name,
		}
	}
	return columns
}

func dtoTasksToResp(dtoTasks []domain.Task) []GetBoardTask {
	tasks := make([]GetBoardTask, len(dtoTasks))
	for i, task := range dtoTasks {
		tasks[i] = GetBoardTask{
			ID:       task.ID,
			ColumnID: task.ColumnID,
			BoardID:  task.BoardID,
			Number:   task.Number,
			Title:    task.Title,
		}
	}
	return tasks
}
