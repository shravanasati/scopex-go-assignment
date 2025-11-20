package service

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkAttendanceRejectsInvalidDate(t *testing.T) {
	rr, resp := performJSONRequest(markAttendance, http.MethodPost, "/attendance", map[string]any{
		"student_id": 1,
		"date":       "2023-13-40",
		"status":     "Present",
	})

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "date must be in YYYY-MM-DD format", resp["error"])
}

func TestGetAttendanceRejectsID(t *testing.T) {
	rr, _ := performJSONRequest(getAttendance, http.MethodGet, "/attendance/should-be-a-number", nil)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
