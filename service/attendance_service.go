package service

import (
	"net/http"
	"strconv"

	model "github.com/shravanasati/scopex-go-assignment/model"
	repository "github.com/shravanasati/scopex-go-assignment/repository"
	util "github.com/shravanasati/scopex-go-assignment/util"

	"github.com/gin-gonic/gin"
)

// RoutesAttendance registers the attendance routes
func RoutesAttendance(rg *gin.RouterGroup) {
	attendance := rg.Group("/attendance")

	attendance.POST("/mark", util.TokenAuthMiddleware(), markAttendance)
	attendance.GET("/:student_id", util.TokenAuthMiddleware(), getAttendance)
}

// markAttendance godoc
// @Summary Mark attendance
// @Description Mark attendance for a student
// @Tags Attendance
// @Accept  json
// @Produce  json
// @Param attendance body model.Attendance true "Attendance"
// @Success 201 {object} model.Attendance
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security bearerAuth
// @Router /attendance/mark [post]
func markAttendance(c *gin.Context) {
	var attendance model.Attendance
	if err := c.ShouldBindJSON(&attendance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate date format if needed, but binding should handle basic string presence.
	// Ideally we parse the date string to ensure it's valid YYYY-MM-DD.

	id, err := repository.MarkAttendance(attendance)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark attendance: " + err.Error()})
		return
	}

	attendance.ID = id
	c.JSON(http.StatusCreated, attendance)
}

// getAttendance godoc
// @Summary Get attendance
// @Description Get attendance records for a student
// @Tags Attendance
// @Accept  json
// @Produce  json
// @Param student_id path int true "Student ID"
// @Success 200 {array} model.Attendance
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security bearerAuth
// @Router /attendance/{student_id} [get]
func getAttendance(c *gin.Context) {
	studentIDStr := c.Param("student_id")
	studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Student ID"})
		return
	}

	attendances, err := repository.GetAttendanceByStudentID(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attendance"})
		return
	}

	c.JSON(http.StatusOK, attendances)
}
