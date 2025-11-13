package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/getboard"
	"github.com/gin-gonic/gin"
)

type (
	GetBoardRequest struct {
		ID string `json:"id"`
	}

	GetBoardResponse struct {
		Name      string          `json:"name"`
		ShortName string          `json:"short_name"`
		Columns   []domain.Column `json:"columns"`
	}

	GetBoardUseCase interface {
		Handle(cmd getboard.GetBoardCommand) (domain.Board, error)
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

	dmn, err := h.getBoardUC.Handle(cmd)
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
		case errors.Is(err, getboard.BoardIsNotExistsErr):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		}
	}

	resp := GetBoardResponse{
		Name:      dmn.Name,
		ShortName: dmn.ShortName,
		Columns:   dmn.Columns,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}
