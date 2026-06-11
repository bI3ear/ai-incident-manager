package handlers

import (
	"ai-incident-manager/database"
	"ai-incident-manager/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const claudeAPI = "https://api.anthropic.com/v1/messages"
const claudeModel = "claude-sonnet-4-6"

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []claudeMessage `json:"messages"`
}

type claudeContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type claudeResponse struct {
	Content []claudeContent `json:"content"`
}

func callClaudeWithMessages(messages []claudeMessage) (string, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	payload := claudeRequest{
		Model:     claudeModel,
		MaxTokens: 1024,
		Messages:  messages,
	}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", claudeAPI, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("claude API error %d: %s", resp.StatusCode, string(respBody))
	}

	var claudeResp claudeResponse
	if err := json.Unmarshal(respBody, &claudeResp); err != nil {
		return "", err
	}

	if len(claudeResp.Content) == 0 {
		return "", fmt.Errorf("empty response from claude")
	}

	return claudeResp.Content[0].Text, nil
}

func callClaude(prompt string) (string, error) {
	return callClaudeWithMessages([]claudeMessage{{Role: "user", Content: prompt}})
}

func AnalyzeIncident(c *gin.Context) {
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

	prompt := fmt.Sprintf(`You are an expert IT incident responder. Analyze this incident and provide a structured response.

Incident Title: %s
Severity: %s
Description: %s

Respond in exactly this format:
## Root Cause Hypothesis
[Your analysis of likely root causes]

## Suggested Fix Steps
1. [Step 1]
2. [Step 2]
3. [Step 3]

## Estimated Resolution Time
[Time estimate with reasoning]`, incident.Title, incident.Severity, incident.Description)

	analysis, err := callClaude(prompt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI analysis failed: " + err.Error()})
		return
	}

	now := time.Now()
	incident.Analysis = analysis
	incident.AnalyzedAt = &now
	if incident.Status == models.StatusOpen {
		incident.Status = models.StatusInProgress
	}

	if err := database.DB.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, incident)
}

func ClassifySeverity(description string) (models.Severity, error) {
	prompt := fmt.Sprintf(`You are an IT operations expert. Classify the severity of this incident.

Description: %s

Severity levels:
- P1: Critical - Complete system outage, data loss, security breach
- P2: High - Major feature broken, significant user impact
- P3: Medium - Minor feature issue, workaround available
- P4: Low - Cosmetic issue, minimal impact

Respond with ONLY one of: P1, P2, P3, P4`, description)

	result, err := callClaude(prompt)
	if err != nil {
		return models.SeverityP3, err
	}

	result = strings.TrimSpace(result)
	switch {
	case strings.Contains(result, "P1"):
		return models.SeverityP1, nil
	case strings.Contains(result, "P2"):
		return models.SeverityP2, nil
	case strings.Contains(result, "P4"):
		return models.SeverityP4, nil
	default:
		return models.SeverityP3, nil
	}
}
