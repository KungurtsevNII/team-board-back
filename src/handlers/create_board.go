package handlers

import (
	"net/http"

	"github.com/KungurtsevNII/team-board-back/src/usecase/createboard"
	"github.com/gin-gonic/gin"
)

type (
	CreateBoardReqest struct {
		Name      string `json:"name"`
		ShortName string `json:"shrort_name"`
	}

	CreateBoardResponce struct {
		ID string `json:"id"`
	}

	CreateBoardUseCase interface {
		Handle(cmd createboard.CreateBoardCommand) (string, error)
	}
)

func (h *HttpHandler) CreateBoard(c *gin.Context) {
	const op = "handlers.CreateBoard"
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
	boardID, err := h.createBoardUC.Handle(*cmd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id": boardID,
	})
}
