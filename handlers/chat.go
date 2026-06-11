package handlers

import (
	"ai-incident-manager/database"
	"ai-incident-manager/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type chatRequest struct {
	Content string `json:"content" binding:"required"`
}

func GetMessages(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var messages []models.Message
	database.DB.Where("incident_id = ?", id).Order("created_at asc").Find(&messages)
	c.JSON(http.StatusOK, messages)
}

func Chat(c *gin.Context) {
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

	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save user message
	userMsg := models.Message{
		IncidentID: uint(id),
		Role:       "user",
		Content:    req.Content,
	}
	if err := database.DB.Create(&userMsg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save message"})
		return
	}

	// Load full conversation history
	var history []models.Message
	database.DB.Where("incident_id = ?", id).Order("created_at asc").Find(&history)

	systemCtx := fmt.Sprintf(
		"You are an expert IT incident responder helping debug and resolve a live incident.\n\nIncident: %s (Severity: %s)\nDescription: %s\n\nYour job is to help the engineer resolve this step by step. Be concise and actionable. When they share new error messages or symptoms, analyze them and suggest next steps.",
		incident.Title, incident.Severity, incident.Description,
	)

	// Build Claude messages: inject system context into the very first user turn
	var claudeMsgs []claudeMessage
	for i, m := range history {
		content := m.Content
		if i == 0 && m.Role == "user" {
			content = systemCtx + "\n\n---\n\n" + content
		}
		claudeMsgs = append(claudeMsgs, claudeMessage{
			Role:    m.Role,
			Content: content,
		})
	}

	reply, err := callClaudeWithMessages(claudeMsgs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI error: " + err.Error()})
		return
	}

	// Save assistant reply
	assistantMsg := models.Message{
		IncidentID: uint(id),
		Role:       "assistant",
		Content:    reply,
	}
	database.DB.Create(&assistantMsg)

	c.JSON(http.StatusOK, gin.H{
		"user":      userMsg,
		"assistant": assistantMsg,
	})
}
