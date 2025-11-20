package service

import (
	"net/http"
	"strconv"
	"strings"

	model "github.com/shravanasati/scopex-go-assignment/model"
	repository "github.com/shravanasati/scopex-go-assignment/repository"
	util "github.com/shravanasati/scopex-go-assignment/util"

	"github.com/gin-gonic/gin"
)

// RoutesStudent registers the student routes
func RoutesStudent(rg *gin.RouterGroup) {
	student := rg.Group("/students")

	student.POST("/", util.TokenAuthMiddleware(), createStudent)
	student.GET("/", util.TokenAuthMiddleware(), getAllStudents)
	student.GET("/:id", util.TokenAuthMiddleware(), getStudentByID)
	student.PUT("/:id", util.TokenAuthMiddleware(), updateStudent)
	student.DELETE("/:id", util.TokenAuthMiddleware(), deleteStudent)
}

// createStudent godoc
// @Summary Create a new student
// @Description Create a new student with the input payload
// @Tags Students
// @Accept  json
// @Produce  json
// @Param student body model.Student true "Student"
// @Success 201 {object} model.Student
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security bearerAuth
// @Router /students/ [post]
func createStudent(c *gin.Context) {
	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := repository.StudentRepo.CreateStudent(student)
	if err != nil {
		// Check for duplicate entry error (MySQL error 1062)
		// This is a simplified check. In a real app, we might check the error code more robustly.
		if strings.Contains(err.Error(), "Duplicate entry") {
			c.JSON(http.StatusConflict, gin.H{"error": "Student already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create student: " + err.Error()})
		return
	}

	student.ID = id
	c.JSON(http.StatusCreated, student)
}

// getAllStudents godoc
// @Summary List all students
// @Description Get a list of students with pagination
// @Tags Students
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Page size" default(10)
// @Success 200 {array} model.Student
// @Failure 500 {object} map[string]string
// @Security bearerAuth
// @Router /students/ [get]
func getAllStudents(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	students, err := repository.StudentRepo.GetAllStudents(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}

	c.JSON(http.StatusOK, students)
}

// getStudentByID godoc
// @Summary Get a student by ID
// @Description Get details of a specific student by ID
// @Tags Students
// @Accept  json
// @Produce  json
// @Param id path int true "Student ID"
// @Success 200 {object} model.Student
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security bearerAuth
// @Router /students/{id} [get]
func getStudentByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	student, err := repository.StudentRepo.GetStudentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	c.JSON(http.StatusOK, student)
}

// updateStudent godoc
// @Summary Update a student
// @Description Update an existing student by ID
// @Tags Students
// @Accept  json
// @Produce  json
// @Param id path int true "Student ID"
// @Param student body model.Student true "Student"
// @Success 200 {object} model.Student
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security bearerAuth
// @Router /students/{id} [put]
func updateStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var student model.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = repository.StudentRepo.UpdateStudent(id, student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	student.ID = id
	c.JSON(http.StatusOK, student)
}

// deleteStudent godoc
// @Summary Delete a student
// @Description Delete a student by ID
// @Tags Students
// @Accept  json
// @Produce  json
// @Param id path int true "Student ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security bearerAuth
// @Router /students/{id} [delete]
func deleteStudent(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	err = repository.StudentRepo.DeleteStudent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}
