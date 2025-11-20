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
	CreateColumnRequest struct {
		Name string `json:"name"`
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

	CreateColumnUseCase interface {
		Handle(
			ctx context.Context,
			cmd createcolumn.Command,
		) (column *domain.Column, err error)
	}
)

// @Summary Создание новой колонки
// @Schemes
// @Tags Columns
// @Accept json
// @Produce json
// @Param board_id path string true "ID доски"
// @Param createColumnRequest body CreateColumnRequest true "request на создание колонки"
// @Success 201 {object}  CreateColumnResponse
// @Failure     400,404,408,500,503  {object}  ErrorResponse
// @Router /v1/boards/{board_id}/columns [POST]
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

	cmd, err := createcolumn.NewCommand(BoardID, req.Name)
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
		case errors.Is(err, createcolumn.ErrBoardIsNotExists):
			NewErrorResponse(c, http.StatusNotFound, "board not found")
		case errors.Is(err, createcolumn.ErrGetLastOrderNumUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to process column order")
		case errors.Is(err, createcolumn.ErrValidationFailed):
			NewErrorResponse(c, http.StatusBadRequest, "validation failed")
		case errors.Is(err, createcolumn.ErrCreateColumnUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, "failed to create column")
		case errors.Is(err, context.Canceled):
			NewErrorResponse(c, http.StatusRequestTimeout, "request canceled")
		case errors.Is(err, context.DeadlineExceeded):
			NewErrorResponse(c, http.StatusServiceUnavailable, "request timeout")
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
