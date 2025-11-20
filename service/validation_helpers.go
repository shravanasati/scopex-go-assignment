package service

import (
	"fmt"
	"strings"
	"time"
)

const isoDateLayout = "2006-01-02"

// validateISODate ensures a date string follows the YYYY-MM-DD layout and
// represents a real calendar date.
func validateISODate(dateStr string) error {
	trimmed := strings.TrimSpace(dateStr)
	if trimmed == "" {
		return fmt.Errorf("date is required")
	}

	if _, err := time.Parse(isoDateLayout, trimmed); err != nil {
		return fmt.Errorf("date must be in YYYY-MM-DD format")
	}

	return nil
}
