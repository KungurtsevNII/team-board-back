package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/usecase/createboard"
	"github.com/gin-gonic/gin"
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
		Handle(cmd createboard.CreateBoardCommand, ctx context.Context) (string, error)
	}
)

func (h *HttpHandler) CreateBoard(c *gin.Context) {
	var req CreateBoardReqest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	cmd, err := createboard.NewCreateBoardCommand(req.Name, req.ShortName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	boardID, err := h.createBoardUC.Handle(cmd, c.Request.Context())
	if err != nil {
		switch {
		case errors.Is(err, createboard.InvalidNameErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, createboard.InvalidShortNameErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, createboard.EmptyNameErr):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, createboard.BoardIsExistsErr):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": boardID,
	})
}
