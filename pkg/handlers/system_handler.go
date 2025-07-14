package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"smg/pkg/models"
	"smg/pkg/services"
)

type SystemHandler struct {
	systemService *services.SystemService
}

func NewSystemHandler(systemService *services.SystemService) *SystemHandler {
	return &SystemHandler{systemService: systemService}
}

func (h *SystemHandler) GetSettings(c *gin.Context) {
	settings, err := h.systemService.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}

func (h *SystemHandler) UpdateSettings(c *gin.Context) {
	var settings map[string]string
	if err := c.ShouldBindJSON(&settings); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.systemService.UpdateSettings(settings)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated successfully"})
}

func (h *SystemHandler) GetStats(c *gin.Context) {
	stats, err := h.systemService.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (h *SystemHandler) GetPlatforms(c *gin.Context) {
	platforms, err := h.systemService.GetPlatforms()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, platforms)
}

func (h *SystemHandler) CreatePlatform(c *gin.Context) {
	var platform models.Platform
	if err := c.ShouldBindJSON(&platform); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdPlatform, err := h.systemService.CreatePlatform(&platform)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdPlatform)
}

func (h *SystemHandler) UpdatePlatform(c *gin.Context) {
	platformID := c.Param("id")
	if platformID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Platform ID is required"})
		return
	}

	var platform models.Platform
	if err := c.ShouldBindJSON(&platform); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedPlatform, err := h.systemService.UpdatePlatform(platformID, &platform)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedPlatform)
}

func (h *SystemHandler) DeletePlatform(c *gin.Context) {
	platformID := c.Param("id")
	if platformID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Platform ID is required"})
		return
	}

	err := h.systemService.DeletePlatform(platformID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Platform deleted successfully"})
}