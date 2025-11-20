package repository

import (
	"database/sql"
	"fmt"
	"log"

	configuration "github.com/shravanasati/scopex-go-assignment/configuration"
	model "github.com/shravanasati/scopex-go-assignment/model"
)

type StudentRepository interface {
	CreateStudent(student model.Student) (int64, error)
	GetAllStudents(limit, offset int) (model.Students, error)
	GetStudentByID(id int64) (model.Student, error)
	UpdateStudent(id int64, student model.Student) error
	DeleteStudent(id int64) error
}

type studentRepository struct{}

var StudentRepo StudentRepository = &studentRepository{}

// CreateStudent inserts a new student into the database
func (r *studentRepository) CreateStudent(student model.Student) (int64, error) {
	db := configuration.DB
	query := "INSERT INTO students (name, email, department) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Println("Error preparing statement: " + err.Error())
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(student.Name, student.Email, student.Department)
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
	var students model.Students

	query := "SELECT id, name, email, department, created_at FROM students LIMIT ? OFFSET ?"
	rows, err := db.Query(query, limit, offset)
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
	var s model.Student

	query := "SELECT id, name, email, department, created_at FROM students WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&s.ID, &s.Name, &s.Email, &s.Department, &s.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return s, fmt.Errorf("student not found")
		}
		log.Println("Error querying student by ID: " + err.Error())
		return s, err
	}

	return s, nil
}

// UpdateStudent updates an existing student
func (r *studentRepository) UpdateStudent(id int64, student model.Student) error {
	db := configuration.DB
	query := "UPDATE students SET name = ?, email = ?, department = ? WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(student.Name, student.Email, student.Department, id)
	if err != nil {
		log.Println("Error updating student: " + err.Error())
		return err
	}

	return nil
}

// DeleteStudent deletes a student by ID
func (r *studentRepository) DeleteStudent(id int64) error {
	db := configuration.DB
	query := "DELETE FROM students WHERE id = ?"
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Println("Error deleting student: " + err.Error())
		return err
	}

	return nil
}
