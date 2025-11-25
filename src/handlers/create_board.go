package handlers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"log/slog"

	"github.com/KungurtsevNII/team-board-back/src/domain"
	"github.com/KungurtsevNII/team-board-back/src/usecase/createboard"
	"github.com/gin-gonic/gin"
)

const(
	nameOfFirstColumn = "TODO"
)
type (
	CreateBoardReqest struct {
		Name      string `json:"name"`
		ShortName string `json:"short_name"`
	}

	CreateBoardResponce struct {
		ID        string               `json:"id"`
		Name      string               `json:"name"`
		ShortName string               `json:"short_name"`
		Ccr       CreateColumnResponse `json:"column"`
		CreatedAt time.Time            `json:"created_at"`
		UpdatedAt time.Time            `json:"updated_at"`
		DeletedAt *time.Time           `json:"deleted_at"`
	}

	CreateBoardUseCase interface {
		Handle(ctx context.Context, cmd createboard.Command) (*domain.Board, *domain.Column, error)
	}
)

// @Summary Создание новой доски
// @Schemes
// @Tags Boards
// @Accept json
// @Produce json
// @Param createBoardRequest body CreateBoardReqest true "request на создание доски"
// @Success 201 {object}  CreateBoardResponce
// @Failure     400,408,409,500,503  {object}  ErrorResponse
// @Router /v1/boards [POST]
func (h *HttpHandler) CreateBoard(c *gin.Context) {
	const op = "handlers.CreateBoard"
	log := slog.Default()
	log.With("op", op)

	var req CreateBoardReqest
	if err := c.BindJSON(&req); err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "invalid request")
		return
	}
	cmd, err := createboard.NewCommand(req.Name, req.ShortName, nameOfFirstColumn)
	if err != nil {
		log.Warn("failed to create command",
			slog.String("err", err.Error()),
			slog.Any("request", req))
		switch {
		case errors.Is(err, createboard.ErrInvalidName):
			NewErrorResponse(c, http.StatusBadRequest, "invalid name")
		case errors.Is(err, createboard.ErrInvalidShortName):
			NewErrorResponse(c, http.StatusBadRequest, "invalid short name")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	board, column, err := h.createBoardUC.Handle(c.Request.Context(), cmd)
	if err != nil {
		log.Warn("failed to create board",
			slog.String("err", err.Error()),
			slog.Any("cmd", cmd))
		switch {
		case errors.Is(err, createboard.ErrValidationFailed):
			NewErrorResponse(c, http.StatusInternalServerError, createboard.ErrValidationFailed.Error())
		case errors.Is(err, createboard.ErrCreateBoard):
			NewErrorResponse(c, http.StatusInternalServerError, createboard.ErrCreateBoard.Error())
		case errors.Is(err, createboard.ErrInvalidColumnName):
			NewErrorResponse(c, http.StatusBadRequest, createboard.ErrInvalidColumnName.Error())
		case errors.Is(err, createboard.ErrGetLastOrderNumUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, createboard.ErrGetLastOrderNumUnknown.Error())
		case errors.Is(err, createboard.ErrCreateColumnUnknown):
			NewErrorResponse(c, http.StatusInternalServerError, createboard.ErrCreateColumnUnknown.Error())
		case errors.Is(err, createboard.ErrInvalidName):
			NewErrorResponse(c, http.StatusBadRequest, createboard.ErrInvalidName.Error())
		case errors.Is(err, createboard.ErrInvalidShortName):
			NewErrorResponse(c, http.StatusBadRequest, createboard.ErrInvalidName.Error())
		case errors.Is(err, createboard.ErrBoardIsExists):
			NewErrorResponse(c, http.StatusConflict, createboard.ErrBoardIsExists.Error())
		case errors.Is(err, context.Canceled):
			NewErrorResponse(c, http.StatusRequestTimeout, "request canceled")
		case errors.Is(err, context.DeadlineExceeded):
			NewErrorResponse(c, http.StatusServiceUnavailable, "request timeout")
		default:
			NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	resp := CreateBoardResponce{
		ID:        board.ID.String(),
		Name:      board.Name,
		ShortName: board.ShortName,
		Ccr: CreateColumnResponse{
			ID:        column.ID.String(),
			BoardID:   column.BoardID.String(),
			Name:      column.Name,
			OrderNum:  column.OrderNum,
			CreatedAt: column.CreatedAt,
			UpdatedAt: column.UpdatedAt,
			DeletedAt: column.DeletedAt,
		},
		CreatedAt: board.CreatedAt,
		UpdatedAt: board.UpdatedAt,
		DeletedAt: board.DeletedAt,
	}

	c.JSON(http.StatusCreated, resp)
}
