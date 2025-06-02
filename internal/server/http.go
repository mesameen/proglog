package server

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mesameen/proglog/internal/model"
)

type Handler struct {
	Log *Log
}

func newHandler() *Handler {
	return &Handler{
		Log: NewLog(),
	}
}

func (h *Handler) handleProduce(c *gin.Context) {
	var req model.ProduceRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Failed to unmarshal the request. Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	off, err := h.Log.Append(req.Record)
	if err != nil {
		log.Printf("Failed to store the request data. Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store log"})
		return
	}
	c.JSON(http.StatusOK, &model.ProduceResponse{
		Offset: off,
	})
}

func (h *Handler) handleConsume(c *gin.Context) {
	var req model.ConsumeRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Failed to unmarshal the request. Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	res, err := h.Log.Read(req.Offset)
	if err != nil && errors.Is(err, model.ErrOffsetNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "requested log doesn't exists"})
		return
	}
	if err != nil {
		log.Printf("Failed to unmarshal the request. Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(http.StatusOK, res)
}
