package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/KungurtsevNII/team-board-back/src/usecase/createcolumn"
)

type (
	// Контракт/Сваггер
	CreateColumnRequest struct {
		Title   string `json:"title"`
		BoardID string  `json:"board_id"`
	}

	CreateColumnResponse struct {
		Title   string `json:"title"`
		BoardID string  `json:"board_id"`
	}

	// Один юз кейс, на один запрос, нра один пользвательский сценарий.
	CreateColumnUseCase interface {
		CreateColumnHandle(cmd createcolumn.CreateColumnCommand) error
	}
)

func (h *HttpHandler) CreateColumn(c *gin.Context) {
	const op = "handlers.Healthcheck"

	log := slog.Default()
	log.With("op", op, "method", c.Request.Method)
	log.Info(c.Request.URL.Path)


	// todo получить параметры из тела
	var req CreateColumnRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	cmd, err := createcolumn.NewCreateColumnCommand(req.Title, req.BoardID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	err = h.columnUC.CreateColumnHandle(cmd)
	if err != nil {
		switch {
		case errors.Is(err, createcolumn.CreateColumnIsExistsErrr):
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

// 1. Handler Request/Response
// 2. Логика приложения. DTO -> use case (работа с базой, работа с кэшом, работа с очередями) -> DTO
// 3.
