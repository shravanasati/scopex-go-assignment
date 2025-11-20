package service

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	model "github.com/shravanasati/scopex-go-assignment/model"
	repository "github.com/shravanasati/scopex-go-assignment/repository"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type studentServiceMock struct {
	mock.Mock
}

func (m *studentServiceMock) CreateStudent(student model.Student) (model.Student, error) {
	args := m.Called(student)
	result, _ := args.Get(0).(model.Student)
	return result, args.Error(1)
}

func (m *studentServiceMock) GetAllStudents(limit, offset int) (model.Students, error) {
	args := m.Called(limit, offset)
	result, _ := args.Get(0).(model.Students)
	return result, args.Error(1)
}

func (m *studentServiceMock) GetStudentByID(id int64) (model.Student, error) {
	args := m.Called(id)
	result, _ := args.Get(0).(model.Student)
	return result, args.Error(1)
}

func (m *studentServiceMock) UpdateStudent(id int64, student model.Student) (model.Student, error) {
	args := m.Called(id, student)
	result, _ := args.Get(0).(model.Student)
	return result, args.Error(1)
}

func (m *studentServiceMock) DeleteStudent(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func init() {
	gin.SetMode(gin.TestMode)
}

func withMockStudentService(t *testing.T, mockSvc StudentService) {
	original := studentSvc
	setStudentService(mockSvc)
	t.Cleanup(func() {
		setStudentService(original)
	})
}

func performJSONRequest(handler gin.HandlerFunc, method, path string, payload any) (*httptest.ResponseRecorder, map[string]any) {
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(method, path, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Request = req

	handler(c)

	var resp map[string]any
	if rr.Body.Len() > 0 {
		_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	}

	return rr, resp
}

func TestCreateStudentReturnsValidationErrors(t *testing.T) {
	mockSvc := &studentServiceMock{}
	validationErr := &ValidationError{Fields: map[string]string{"email": "invalid"}}
	mockSvc.On("CreateStudent", mock.AnythingOfType("model.Student")).Return(model.Student{}, validationErr)

	withMockStudentService(t, mockSvc)

	rr, resp := performJSONRequest(createStudent, http.MethodPost, "/students", map[string]any{
		"name":       "Jane",
		"email":      "jane@example.com",
		"department": "Science",
	})

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Equal(t, "validation failed", resp["error"])
	mockSvc.AssertExpectations(t)
}

func TestCreateStudentRejectsDuplicateEmails(t *testing.T) {
	mockSvc := &studentServiceMock{}
	mockSvc.On("CreateStudent", mock.AnythingOfType("model.Student")).Return(model.Student{}, ErrDuplicateStudentEmail)

	withMockStudentService(t, mockSvc)

	rr, resp := performJSONRequest(createStudent, http.MethodPost, "/students", map[string]any{
		"name":       "Jane",
		"email":      "jane@example.com",
		"department": "Science",
	})

	assert.Equal(t, http.StatusConflict, rr.Code)
	assert.Equal(t, ErrDuplicateStudentEmail.Error(), resp["error"])
	mockSvc.AssertExpectations(t)
}

func TestCreateStudentSuccess(t *testing.T) {
	mockSvc := &studentServiceMock{}
	expected := model.Student{ID: 1, Name: "Jane", Email: "jane@example.com", Department: "Science"}
	mockSvc.On("CreateStudent", mock.AnythingOfType("model.Student")).Return(expected, nil)

	withMockStudentService(t, mockSvc)

	rr, resp := performJSONRequest(createStudent, http.MethodPost, "/students", map[string]any{
		"name":       "Jane",
		"email":      "jane@example.com",
		"department": "Science",
	})

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Equal(t, float64(1), resp["id"])
	mockSvc.AssertExpectations(t)
}

func TestDeleteStudentReturnsNotFound(t *testing.T) {
	mockSvc := &studentServiceMock{}
	mockSvc.On("DeleteStudent", int64(9)).Return(repository.ErrStudentNotFound)

	withMockStudentService(t, mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/students/9", nil)
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "9"}}
	c.Request = req

	deleteStudent(c)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockSvc.AssertExpectations(t)
}
