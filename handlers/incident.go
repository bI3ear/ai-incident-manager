package handlers

import (
	"ai-incident-manager/database"
	"ai-incident-manager/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ListIncidents(c *gin.Context) {
	var incidents []models.Incident
	if err := database.DB.Order("created_at desc").Find(&incidents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

func GetIncident(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var incident models.Incident
	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
		return
	}
	c.JSON(http.StatusOK, incident)
}

func CreateIncident(c *gin.Context) {
	var req models.CreateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Use AI to suggest severity
	severity, err := ClassifySeverity(req.Description)
	if err != nil {
		severity = models.SeverityP3 // fallback
	}

	incident := models.Incident{
		Title:       req.Title,
		Description: req.Description,
		Severity:    severity,
		Status:      models.StatusOpen,
	}

	if err := database.DB.Create(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, incident)
}

func UpdateIncident(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var incident models.Incident
	if err := database.DB.First(&incident, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "incident not found"})
		return
	}

	var req models.UpdateIncidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Title != "" {
		incident.Title = req.Title
	}
	if req.Description != "" {
		incident.Description = req.Description
	}
	if req.Severity != "" {
		incident.Severity = req.Severity
	}
	if req.Status != "" {
		incident.Status = req.Status
	}

	if err := database.DB.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, incident)
}

func DeleteIncident(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := database.DB.Delete(&models.Incident{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "incident deleted"})
}
