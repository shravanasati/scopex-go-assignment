package service

import (
	"errors"
	"net/mail"
	"strings"

	model "github.com/shravanasati/scopex-go-assignment/model"
	repository "github.com/shravanasati/scopex-go-assignment/repository"
)

// ErrDuplicateStudentEmail is returned when attempting to store a student
// record with an email address that already exists.
var ErrDuplicateStudentEmail = errors.New("student with this email already exists")

// StudentService describes the student domain operations that the HTTP layer
// relies on. Having an explicit interface allows us to unit-test handlers by
// supplying mocks.
type StudentService interface {
	CreateStudent(student model.Student) (model.Student, error)
	GetAllStudents(limit, offset int) (model.Students, error)
	GetStudentByID(id int64) (model.Student, error)
	UpdateStudent(id int64, student model.Student) (model.Student, error)
	DeleteStudent(id int64) error
}

type studentService struct {
	repo repository.StudentRepository
}

var studentSvc StudentService = newStudentService(repository.StudentRepo)

func newStudentService(repo repository.StudentRepository) StudentService {
	return &studentService{repo: repo}
}

// setStudentService allows tests to inject a mock implementation.
func setStudentService(svc StudentService) {
	studentSvc = svc
}

func (s *studentService) CreateStudent(student model.Student) (model.Student, error) {
	if err := validateStudentInput(student); err != nil {
		return model.Student{}, err
	}

	existing, err := s.repo.GetStudentByEmail(student.Email)
	if err != nil {
		return model.Student{}, err
	}
	if existing.ID != 0 {
		return model.Student{}, ErrDuplicateStudentEmail
	}

	id, err := s.repo.CreateStudent(student)
	if err != nil {
		return model.Student{}, err
	}

	student.ID = id
	return student, nil
}

func (s *studentService) GetAllStudents(limit, offset int) (model.Students, error) {
	return s.repo.GetAllStudents(limit, offset)
}

func (s *studentService) GetStudentByID(id int64) (model.Student, error) {
	return s.repo.GetStudentByID(id)
}

func (s *studentService) UpdateStudent(id int64, student model.Student) (model.Student, error) {
	if err := validateStudentInput(student); err != nil {
		return model.Student{}, err
	}

	existing, err := s.repo.GetStudentByEmail(student.Email)
	if err != nil {
		return model.Student{}, err
	}
	if existing.ID != 0 && existing.ID != id {
		return model.Student{}, ErrDuplicateStudentEmail
	}

	if err := s.repo.UpdateStudent(id, student); err != nil {
		return model.Student{}, err
	}

	student.ID = id
	return student, nil
}

func (s *studentService) DeleteStudent(id int64) error {
	return s.repo.DeleteStudent(id)
}

// ValidationError captures field-level validation failures.
type ValidationError struct {
	Fields map[string]string
}

func (v *ValidationError) Error() string {
	return "validation failed"
}

func validateStudentInput(student model.Student) error {
	issues := make(map[string]string)

	if strings.TrimSpace(student.Name) == "" {
		issues["name"] = "name is required"
	}
	if strings.TrimSpace(student.Email) == "" {
		issues["email"] = "email is required"
	} else if _, err := mail.ParseAddress(student.Email); err != nil {
		issues["email"] = "email format is invalid"
	}
	if strings.TrimSpace(student.Department) == "" {
		issues["department"] = "department is required"
	}

	if len(issues) > 0 {
		return &ValidationError{Fields: issues}
	}

	return nil
}
