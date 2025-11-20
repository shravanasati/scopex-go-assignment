package repository

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	configuration "github.com/shravanasati/scopex-go-assignment/configuration"
	model "github.com/shravanasati/scopex-go-assignment/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupStudentSQLMock(t *testing.T) (sqlmock.Sqlmock, *sql.DB) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	configuration.DB = db
	t.Cleanup(func() {
		db.Close()
	})

	return mock, db
}

func TestGetStudentByEmailFound(t *testing.T) {
	mock, _ := setupStudentSQLMock(t)
	repo := &studentRepository{}

	email := "jane@example.com"
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "email", "department", "created_at"}).
		AddRow(int64(1), "Jane", email, "Science", now)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, department, created_at FROM students WHERE email = ?")).
		WithArgs(email).
		WillReturnRows(rows)

	student, err := repo.GetStudentByEmail(email)

	assert.NoError(t, err)
	assert.Equal(t, email, student.Email)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetStudentByEmailNotFound(t *testing.T) {
	mock, _ := setupStudentSQLMock(t)
	repo := &studentRepository{}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, name, email, department, created_at FROM students WHERE email = ?")).
		WithArgs("ghost@example.com").
		WillReturnError(sql.ErrNoRows)

	student, err := repo.GetStudentByEmail("ghost@example.com")

	assert.NoError(t, err)
	assert.Equal(t, int64(0), student.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateStudentNotFound(t *testing.T) {
	mock, _ := setupStudentSQLMock(t)
	repo := &studentRepository{}

	input := model.Student{Name: "Jane", Email: "jane@example.com", Department: "Science"}

	prep := mock.ExpectPrepare(regexp.QuoteMeta("UPDATE students SET name = ?, email = ?, department = ? WHERE id = ?"))
	prep.ExpectExec().
		WithArgs(input.Name, input.Email, input.Department, int64(99)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.UpdateStudent(99, input)

	assert.ErrorIs(t, err, ErrStudentNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteStudentNotFound(t *testing.T) {
	mock, _ := setupStudentSQLMock(t)
	repo := &studentRepository{}

	prep := mock.ExpectPrepare(regexp.QuoteMeta("DELETE FROM students WHERE id = ?"))
	prep.ExpectExec().
		WithArgs(int64(77)).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err := repo.DeleteStudent(77)

	assert.ErrorIs(t, err, ErrStudentNotFound)
	assert.NoError(t, mock.ExpectationsWereMet())
}
