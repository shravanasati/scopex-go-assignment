package service

import (
	"errors"
	"testing"

	model "github.com/shravanasati/scopex-go-assignment/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockStudentRepository struct {
	mock.Mock
}

func (m *mockStudentRepository) CreateStudent(student model.Student) (int64, error) {
	args := m.Called(student)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockStudentRepository) GetAllStudents(limit, offset int) (model.Students, error) {
	args := m.Called(limit, offset)
	if students, ok := args.Get(0).(model.Students); ok {
		return students, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockStudentRepository) GetStudentByID(id int64) (model.Student, error) {
	args := m.Called(id)
	return args.Get(0).(model.Student), args.Error(1)
}

func (m *mockStudentRepository) GetStudentByEmail(email string) (model.Student, error) {
	args := m.Called(email)
	return args.Get(0).(model.Student), args.Error(1)
}

func (m *mockStudentRepository) UpdateStudent(id int64, student model.Student) error {
	args := m.Called(id, student)
	return args.Error(0)
}

func (m *mockStudentRepository) DeleteStudent(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestStudentServiceCreateStudentSuccess(t *testing.T) {
	repo := &mockStudentRepository{}
	svc := newStudentService(repo)

	input := model.Student{Name: "Jane", Email: "jane@example.com", Department: "Science"}

	repo.On("GetStudentByEmail", input.Email).Return(model.Student{}, nil).Once()
	repo.On("CreateStudent", input).Return(int64(42), nil).Once()

	created, err := svc.CreateStudent(input)

	assert.NoError(t, err)
	assert.Equal(t, int64(42), created.ID)
	repo.AssertExpectations(t)
}

func TestStudentServiceCreateStudentInvalidEmail(t *testing.T) {
	repo := &mockStudentRepository{}
	svc := newStudentService(repo)

	input := model.Student{Name: "Jane", Email: "invalid-email", Department: "Science"}

	_, err := svc.CreateStudent(input)

	var validationErr *ValidationError
	assert.True(t, errors.As(err, &validationErr))
	repo.AssertExpectations(t)
}

func TestStudentServiceCreateStudentDuplicateEmail(t *testing.T) {
	repo := &mockStudentRepository{}
	svc := newStudentService(repo)

	input := model.Student{Name: "Jane", Email: "jane@example.com", Department: "Science"}

	repo.On("GetStudentByEmail", input.Email).Return(model.Student{ID: 7, Email: input.Email}, nil).Once()

	_, err := svc.CreateStudent(input)

	assert.ErrorIs(t, err, ErrDuplicateStudentEmail)
	repo.AssertExpectations(t)
}

func TestStudentServiceUpdateStudentDuplicateEmail(t *testing.T) {
	repo := &mockStudentRepository{}
	svc := newStudentService(repo)

	input := model.Student{Name: "Jane", Email: "jane@example.com", Department: "Science"}

	repo.On("GetStudentByEmail", input.Email).Return(model.Student{ID: 99, Email: input.Email}, nil).Once()

	_, err := svc.UpdateStudent(1, input)

	assert.ErrorIs(t, err, ErrDuplicateStudentEmail)
	repo.AssertExpectations(t)
}

func TestStudentServiceUpdateStudentSuccess(t *testing.T) {
	repo := &mockStudentRepository{}
	svc := newStudentService(repo)

	input := model.Student{Name: "Jane", Email: "jane@example.com", Department: "Science"}

	repo.On("GetStudentByEmail", input.Email).Return(model.Student{ID: 1, Email: input.Email}, nil).Once()
	repo.On("UpdateStudent", int64(1), input).Return(nil).Once()

	updated, err := svc.UpdateStudent(1, input)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), updated.ID)
	repo.AssertExpectations(t)
}
