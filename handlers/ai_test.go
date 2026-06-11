package handlers

import (
	"testing"

	"ai-incident-manager/models"
)

func TestClassifySeverityFallback(t *testing.T) {
	// Without API key, should return fallback P3
	sev, _ := ClassifySeverity("some description")
	if sev == "" {
		t.Error("expected a severity value, got empty string")
	}
}

func TestSeverityConstants(t *testing.T) {
	cases := []models.Severity{
		models.SeverityP1,
		models.SeverityP2,
		models.SeverityP3,
		models.SeverityP4,
	}
	for _, s := range cases {
		if s == "" {
			t.Errorf("severity constant should not be empty")
		}
	}
}
