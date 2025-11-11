package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/getcolumn"
)

type (
	// Контракт/Сваггер

	GetColumnResponse struct {
		Title   string `json:"title"`
		BoardID string  `json:"board_id"`
	}

	// Один юз кейс, на один запрос, нра один пользвательский сценарий.
	GetColumnUseCase interface {
		Handle(cmd getcolumn.GetColumnCommand) (domain.Column,error)
	}
)

func (h *HttpHandler) GetColumn(c *gin.Context) {
	const op = "handlers.Healthcheck"

	log := slog.Default()
	log.With("op", op, "method", c.Request.Method)
	log.Info(c.Request.URL.Path)

	// todo получить параметры из тела
	id := c.Param("id")

	cmd, err := getcolumn.NewGetColumnCommand(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	dmn, err := h.getColumnUC.Handle(cmd)
	if err != nil {
		switch {
		// case errors.Is(err, createcolumn.ErrColumnIsExistsErr):
		// 	c.JSON(http.StatusConflict, gin.H{
		// 		"error": err.Error(),
		// 	})
		}
	}

	resp := GetColumnResponse{
		Title:   dmn.Name,
		// BoardID: dmn.BoardID,
	}

	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"data":   resp,
	})
}

// 1. Handler Request/Response
// 2. Логика приложения. DTO -> use case (работа с базой, работа с кэшом, работа с очередями) -> DTO
// 3.
