package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/usecase/createboard"
	"github.com/gin-gonic/gin"
	"log/slog"
)

type (
	CreateBoardReqest struct {
		Name      string `json:"name"`
		ShortName string `json:"short_name"`
	}

	CreateBoardResponce struct {
		ID string `json:"id"`
	}

	CreateBoardUseCase interface {
		Handle(ctx context.Context, cmd createboard.Command) (string, error)
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
	cmd, err := createboard.NewCommand(req.Name, req.ShortName)
	if err != nil {
		NewErrorResponse(c, http.StatusBadRequest, "failed to create command")
		return
	}

	boardID, err := h.createBoardUC.Handle(c.Request.Context(), cmd)
	if err != nil {
		switch {
		case errors.Is(err, createboard.ErrInvalidName):
			NewErrorResponse(c, http.StatusBadRequest, createboard.ErrInvalidName.Error())
		case errors.Is(err, createboard.ErrInvalidShortName):
			NewErrorResponse(c, http.StatusBadRequest, createboard.ErrInvalidName.Error())
		case errors.Is(err, createboard.ErrEmptyName):
			NewErrorResponse(c, http.StatusBadRequest, createboard.ErrEmptyName.Error())
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

	c.JSON(http.StatusCreated, CreateBoardResponce{ID: boardID})
}
