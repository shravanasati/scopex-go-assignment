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

func TestStudentInvalidID(t *testing.T) {
	rr, _ := performJSONRequest(updateStudent, http.MethodPut, "/students/invalid-id", nil)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	rr, _ = performJSONRequest(updateStudent, http.MethodPut, "/students/8", map[string]string{})
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	rr, _ = performJSONRequest(deleteStudent, http.MethodDelete, "/students/invalid-id", nil)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	rr, _ = performJSONRequest(getStudentByID, http.MethodGet, "/students/invalid-id", nil)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	rr, _ = performJSONRequest(createStudent, http.MethodPost, "/students/invalid-id", nil)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestGetAllStudentsSuccess(t *testing.T) {
	mockSvc := &studentServiceMock{}
	expected := model.Students{
		{ID: 1, Name: "John", Email: "john@example.com", Department: "Math"},
		{ID: 2, Name: "Jane", Email: "jane@example.com", Department: "Science"},
	}
	mockSvc.On("GetAllStudents", 10, 0).Return(expected, nil)

	withMockStudentService(t, mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/students?page=1&limit=10", nil)
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Request = req

	getAllStudents(c)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp model.Students
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.Len(t, resp, 2)
	mockSvc.AssertExpectations(t)
}

func TestGetStudentByIDSuccess(t *testing.T) {
	mockSvc := &studentServiceMock{}
	expected := model.Student{ID: 1, Name: "John", Email: "john@example.com", Department: "Math"}
	mockSvc.On("GetStudentByID", int64(1)).Return(expected, nil)

	withMockStudentService(t, mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/students/1", nil)
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	getStudentByID(c)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp model.Student
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.Equal(t, int64(1), resp.ID)
	mockSvc.AssertExpectations(t)
}

func TestGetStudentByIDNotFound(t *testing.T) {
	mockSvc := &studentServiceMock{}
	mockSvc.On("GetStudentByID", int64(999)).Return(model.Student{}, repository.ErrStudentNotFound)

	withMockStudentService(t, mockSvc)

	req := httptest.NewRequest(http.MethodGet, "/students/999", nil)
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "999"}}
	c.Request = req

	getStudentByID(c)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdateStudentSuccess(t *testing.T) {
	mockSvc := &studentServiceMock{}
	expected := model.Student{ID: 1, Name: "John Updated", Email: "john@example.com", Department: "Math"}
	mockSvc.On("UpdateStudent", int64(1), mock.AnythingOfType("model.Student")).Return(expected, nil)

	withMockStudentService(t, mockSvc)

	payload := map[string]any{
		"name":       "John Updated",
		"email":      "john@example.com",
		"department": "Math",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/students/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	updateStudent(c)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp model.Student
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.Equal(t, "John Updated", resp.Name)
	mockSvc.AssertExpectations(t)
}

func TestUpdateStudentValidationErrors(t *testing.T) {
	mockSvc := &studentServiceMock{}
	validationErr := &ValidationError{Fields: map[string]string{"email": "invalid"}}
	mockSvc.On("UpdateStudent", int64(1), mock.AnythingOfType("model.Student")).Return(model.Student{}, validationErr)

	withMockStudentService(t, mockSvc)

	payload := map[string]any{
		"name":       "John",
		"email":      "john@example.com",
		"department": "Math",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/students/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	updateStudent(c)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.Equal(t, "validation failed", resp["error"])
	mockSvc.AssertExpectations(t)
}

func TestUpdateStudentNotFound(t *testing.T) {
	mockSvc := &studentServiceMock{}
	mockSvc.On("UpdateStudent", int64(999), mock.AnythingOfType("model.Student")).Return(model.Student{}, repository.ErrStudentNotFound)

	withMockStudentService(t, mockSvc)

	payload := map[string]any{
		"name":       "John",
		"email":      "john@example.com",
		"department": "Math",
	}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPut, "/students/999", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "999"}}
	c.Request = req

	updateStudent(c)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.Equal(t, repository.ErrStudentNotFound.Error(), resp["error"])
	mockSvc.AssertExpectations(t)
}

func TestDeleteStudentSuccess(t *testing.T) {
	mockSvc := &studentServiceMock{}
	mockSvc.On("DeleteStudent", int64(1)).Return(nil)

	withMockStudentService(t, mockSvc)

	req := httptest.NewRequest(http.MethodDelete, "/students/1", nil)
	rr := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rr)
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}
	c.Request = req

	deleteStudent(c)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp map[string]any
	_ = json.Unmarshal(rr.Body.Bytes(), &resp)
	assert.Equal(t, "Student deleted successfully", resp["message"])
	mockSvc.AssertExpectations(t)
}
