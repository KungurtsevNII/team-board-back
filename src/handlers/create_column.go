package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/createcolumn"
)

type (
	// Контракт/Сваггер
	CreateColumnRequest struct {
		Name     string `json:"name"`
	}

	CreateColumnResponse struct {
		ID        string     `json:"id"`
		BoardID   string     `json:"board_id"`
		OrderNum  int64      `json:"order_num"`
		Name      string     `json:"name"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at"`
	}

	// Один юз кейс, на один запрос, нра один пользвательский сценарий.
	CreateColumnUseCase interface {
		Handle(
			ctx context.Context,
			cmd createcolumn.CreateColumnCommand,
		) (column *domain.Column, err error)
	}
)

func (h *HttpHandler) CreateColumn(c *gin.Context) {
	const op = "handlers.CreateColumn"
	log := slog.Default()
	log.With("op", op)

	BoardID := c.Param("board_id")

	var req CreateColumnRequest
	if err := c.BindJSON(&req); err != nil {
		log.Warn("failed to bind request", slog.String("err", err.Error()))
		NewErrorResponse(c, http.StatusBadRequest, "bad body")
		return
	}

	cmd, err := createcolumn.NewCreateColumnCommand(BoardID, req.Name)
	if err != nil {
		log.Warn("failed to create command",
			slog.String("err", err.Error()),
			slog.String("board_id", BoardID),
			slog.String("name", req.Name))
		
		switch {
		case errors.Is(err, createcolumn.ErrValidationFailed):
			NewErrorResponse(c, http.StatusBadRequest, "validation failed")
		case errors.Is(err, createcolumn.ErrInvalidUUID):
			NewErrorResponse(c, http.StatusBadRequest, "invalid board id")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	dmn, err := h.createColumnUC.Handle(c, cmd)
	if err != nil {
		log.Error("failed to create column",
			slog.String("err", err.Error()),
			slog.String("board_id", cmd.BoardID.String()),
			slog.String("name", cmd.Name))
		
		switch {
		case errors.Is(err, createcolumn.ErrBoardIsNotExistsErr):
			NewErrorResponse(c, http.StatusNotFound, "board not found")
		case errors.Is(err, createcolumn.ErrGetLastOrderNumErr):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to process column order")
		case errors.Is(err, createcolumn.ErrValidationFailed):
			NewErrorResponse(c, http.StatusBadRequest, "validation failed")
		case errors.Is(err, createcolumn.ErrCreateColumnErr):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to create column")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	resp := CreateColumnResponse{
		ID:        dmn.ID.String(),
		BoardID:   dmn.BoardID.String(),
		OrderNum:  dmn.OrderNum,
		Name:      dmn.Name,
		CreatedAt: dmn.CreatedAt,
		UpdatedAt: dmn.UpdatedAt,
		DeletedAt: dmn.DeletedAt,
	}

	c.JSON(http.StatusCreated, resp)
}

// 1. Handler Request/Response
// 2. Логика приложения. DTO -> use case (работа с базой, работа с кэшом, работа с очередями) -> DTO
// 3.
