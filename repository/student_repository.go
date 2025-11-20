package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	configuration "github.com/shravanasati/scopex-go-assignment/configuration"
	model "github.com/shravanasati/scopex-go-assignment/model"
)

type StudentRepository interface {
	CreateStudent(student model.Student) (int64, error)
	GetAllStudents(limit, offset int) (model.Students, error)
	GetStudentByID(id int64) (model.Student, error)
	GetStudentByEmail(email string) (model.Student, error)
	UpdateStudent(id int64, student model.Student) error
	DeleteStudent(id int64) error
}
type studentRepository struct{}

var StudentRepo StudentRepository = &studentRepository{}

// ErrStudentNotFound indicates that the requested student record does not
// exist in persistent storage.
var ErrStudentNotFound = errors.New("student not found")

// CreateStudent inserts a new student into the database
func (r *studentRepository) CreateStudent(student model.Student) (int64, error) {
	db := configuration.DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO students (name, email, department) VALUES (?, ?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Println("Error preparing statement: " + err.Error())
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, student.Name, student.Email, student.Department)
	if err != nil {
		log.Println("Error inserting student: " + err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetAllStudents retrieves all students with pagination
func (r *studentRepository) GetAllStudents(limit, offset int) (model.Students, error) {
	db := configuration.DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var students model.Students

	query := "SELECT id, name, email, department, created_at FROM students LIMIT ? OFFSET ?"
	rows, err := db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		log.Println("Error querying students: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s model.Student
		// Scan created_at as []uint8 (byte slice) if driver returns it as such, or time.Time if configured.
		// The mysql driver usually handles time.Time if parseTime=true is in DSN.
		// Let's assume standard scanning works.
		err := rows.Scan(&s.ID, &s.Name, &s.Email, &s.Department, &s.CreatedAt)
		if err != nil {
			log.Println("Error scanning student: " + err.Error())
			return nil, err
		}
		students = append(students, s)
	}

	return students, nil
}

// GetStudentByID retrieves a student by ID
func (r *studentRepository) GetStudentByID(id int64) (model.Student, error) {
	db := configuration.DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var s model.Student

	query := "SELECT id, name, email, department, created_at FROM students WHERE id = ?"
	err := db.QueryRowContext(ctx, query, id).Scan(&s.ID, &s.Name, &s.Email, &s.Department, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return s, ErrStudentNotFound
		}
		log.Println("Error querying student by ID: " + err.Error())
		return s, err
	}

	return s, nil
}

// GetStudentByEmail retrieves a student record by email address.
func (r *studentRepository) GetStudentByEmail(email string) (model.Student, error) {
	db := configuration.DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var s model.Student

	query := "SELECT id, name, email, department, created_at FROM students WHERE email = ?"
	err := db.QueryRowContext(ctx, query, email).Scan(&s.ID, &s.Name, &s.Email, &s.Department, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Student{}, nil
		}
		return model.Student{}, err
	}

	return s, nil
}

// UpdateStudent updates an existing student
func (r *studentRepository) UpdateStudent(id int64, student model.Student) error {
	db := configuration.DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "UPDATE students SET name = ?, email = ?, department = ? WHERE id = ?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, student.Name, student.Email, student.Department, id)
	if err != nil {
		log.Println("Error updating student: " + err.Error())
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrStudentNotFound
	}

	return nil
}

// DeleteStudent deletes a student by ID
func (r *studentRepository) DeleteStudent(id int64) error {
	db := configuration.DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "DELETE FROM students WHERE id = ?"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		log.Println("Error deleting student: " + err.Error())
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return ErrStudentNotFound
	}

	return nil
}
