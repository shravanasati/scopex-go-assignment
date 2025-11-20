package repository

import (
	"database/sql"
	"regexp"
	"testing"

	configuration "github.com/shravanasati/scopex-go-assignment/configuration"
	model "github.com/shravanasati/scopex-go-assignment/model"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func setupAttendanceSQLMock(t *testing.T) (sqlmock.Sqlmock, *sql.DB) {
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

func TestMarkAttendanceSuccess(t *testing.T) {
	mock, _ := setupAttendanceSQLMock(t)

	attendance := model.Attendance{
		StudentID: 1,
		Date:      "2023-10-27",
		Status:    "Present",
	}

	prep := mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO attendance (student_id, date, status) VALUES (?, ?, ?)"))
	prep.ExpectExec().
		WithArgs(attendance.StudentID, attendance.Date, attendance.Status).
		WillReturnResult(sqlmock.NewResult(1, 1))

	id, err := MarkAttendance(attendance)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestMarkAttendanceError(t *testing.T) {
	mock, _ := setupAttendanceSQLMock(t)

	attendance := model.Attendance{
		StudentID: 1,
		Date:      "2023-10-27",
		Status:    "Present",
	}

	prep := mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO attendance (student_id, date, status) VALUES (?, ?, ?)"))
	prep.ExpectExec().
		WithArgs(attendance.StudentID, attendance.Date, attendance.Status).
		WillReturnError(sql.ErrConnDone)

	id, err := MarkAttendance(attendance)

	assert.Error(t, err)
	assert.Equal(t, int64(0), id)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAttendanceByStudentIDSuccess(t *testing.T) {
	mock, _ := setupAttendanceSQLMock(t)

	studentID := int64(1)
	rows := sqlmock.NewRows([]string{"id", "student_id", "date", "status"}).
		AddRow(int64(1), studentID, "2023-10-27", "Present").
		AddRow(int64(2), studentID, "2023-10-26", "Absent")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, student_id, date, status FROM attendance WHERE student_id = ? ORDER BY date DESC")).
		WithArgs(studentID).
		WillReturnRows(rows)

	attendances, err := GetAttendanceByStudentID(studentID)

	assert.NoError(t, err)
	assert.Len(t, attendances, 2)
	assert.Equal(t, "2023-10-27", attendances[0].Date)
	assert.Equal(t, "Present", attendances[0].Status)
	assert.Equal(t, "Absent", attendances[1].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAttendanceByStudentIDNotFound(t *testing.T) {
	mock, _ := setupAttendanceSQLMock(t)

	studentID := int64(999)
	rows := sqlmock.NewRows([]string{"id", "student_id", "date", "status"})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, student_id, date, status FROM attendance WHERE student_id = ? ORDER BY date DESC")).
		WithArgs(studentID).
		WillReturnRows(rows)

	attendances, err := GetAttendanceByStudentID(studentID)

	assert.NoError(t, err)
	assert.Len(t, attendances, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAttendanceByStudentIDError(t *testing.T) {
	mock, _ := setupAttendanceSQLMock(t)

	studentID := int64(1)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, student_id, date, status FROM attendance WHERE student_id = ? ORDER BY date DESC")).
		WithArgs(studentID).
		WillReturnError(sql.ErrConnDone)

	attendances, err := GetAttendanceByStudentID(studentID)

	assert.Error(t, err)
	assert.Nil(t, attendances)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAttendanceByDateRangeSuccess(t *testing.T) {
	mock, _ := setupAttendanceSQLMock(t)

	studentID := int64(1)
	startDate := "2023-10-01"
	endDate := "2023-10-31"

	rows := sqlmock.NewRows([]string{"id", "student_id", "date", "status"}).
		AddRow(int64(1), studentID, "2023-10-10", "Present").
		AddRow(int64(2), studentID, "2023-10-15", "Absent").
		AddRow(int64(3), studentID, "2023-10-20", "Present")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, student_id, date, status FROM attendance WHERE student_id = ? AND date BETWEEN ? AND ? ORDER BY date ASC")).
		WithArgs(studentID, startDate, endDate).
		WillReturnRows(rows)

	attendances, err := GetAttendanceByDateRange(studentID, startDate, endDate)

	assert.NoError(t, err)
	assert.Len(t, attendances, 3)
	assert.Equal(t, "2023-10-10", attendances[0].Date)
	assert.Equal(t, "Present", attendances[0].Status)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAttendanceByDateRangeEmpty(t *testing.T) {
	mock, _ := setupAttendanceSQLMock(t)

	studentID := int64(1)
	startDate := "2023-01-01"
	endDate := "2023-01-31"

	rows := sqlmock.NewRows([]string{"id", "student_id", "date", "status"})

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, student_id, date, status FROM attendance WHERE student_id = ? AND date BETWEEN ? AND ? ORDER BY date ASC")).
		WithArgs(studentID, startDate, endDate).
		WillReturnRows(rows)

	attendances, err := GetAttendanceByDateRange(studentID, startDate, endDate)

	assert.NoError(t, err)
	assert.Len(t, attendances, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAttendanceByDateRangeError(t *testing.T) {
	mock, _ := setupAttendanceSQLMock(t)

	studentID := int64(1)
	startDate := "2023-10-01"
	endDate := "2023-10-31"

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, student_id, date, status FROM attendance WHERE student_id = ? AND date BETWEEN ? AND ? ORDER BY date ASC")).
		WithArgs(studentID, startDate, endDate).
		WillReturnError(sql.ErrConnDone)

	attendances, err := GetAttendanceByDateRange(studentID, startDate, endDate)

	assert.Error(t, err)
	assert.Nil(t, attendances)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAttendanceReportSuccess(t *testing.T) {
	mock, _ := setupAttendanceSQLMock(t)

	startDate := "2023-10-01"
	endDate := "2023-10-31"

	rows := sqlmock.NewRows([]string{"id", "name", "email", "present_count", "absent_count"}).
		AddRow(int64(1), "Alice Smith", "alice@example.com", 15, 5).
		AddRow(int64(2), "Bob Johnson", "bob@example.com", 18, 2).
		AddRow(int64(3), "Charlie Brown", "charlie@example.com", 10, 10)

	expectedQuery := `
		SELECT 
			s.id, 
			s.name, 
			s.email, 
			COALESCE(SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END), 0) as present_count,
			COALESCE(SUM(CASE WHEN a.status = 'Absent' THEN 1 ELSE 0 END), 0) as absent_count
		FROM 
			students s
		LEFT JOIN 
			attendance a ON s.id = a.student_id AND a.date BETWEEN ? AND ?
		GROUP BY 
			s.id, s.name, s.email
		ORDER BY 
			s.name ASC
	`

	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(startDate, endDate).
		WillReturnRows(rows)

	reports, err := GetAttendanceReport(startDate, endDate)

	assert.NoError(t, err)
	assert.Len(t, reports, 3)
	assert.Equal(t, "Alice Smith", reports[0].StudentName)
	assert.Equal(t, 15, reports[0].PresentCount)
	assert.Equal(t, 5, reports[0].AbsentCount)
	assert.Equal(t, "Bob Johnson", reports[1].StudentName)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAttendanceReportEmpty(t *testing.T) {
	mock, _ := setupAttendanceSQLMock(t)

	startDate := "2023-10-01"
	endDate := "2023-10-31"

	rows := sqlmock.NewRows([]string{"id", "name", "email", "present_count", "absent_count"})

	expectedQuery := `
		SELECT 
			s.id, 
			s.name, 
			s.email, 
			COALESCE(SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END), 0) as present_count,
			COALESCE(SUM(CASE WHEN a.status = 'Absent' THEN 1 ELSE 0 END), 0) as absent_count
		FROM 
			students s
		LEFT JOIN 
			attendance a ON s.id = a.student_id AND a.date BETWEEN ? AND ?
		GROUP BY 
			s.id, s.name, s.email
		ORDER BY 
			s.name ASC
	`

	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(startDate, endDate).
		WillReturnRows(rows)

	reports, err := GetAttendanceReport(startDate, endDate)

	assert.NoError(t, err)
	assert.Len(t, reports, 0)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAttendanceReportError(t *testing.T) {
	mock, _ := setupAttendanceSQLMock(t)

	startDate := "2023-10-01"
	endDate := "2023-10-31"

	expectedQuery := `
		SELECT 
			s.id, 
			s.name, 
			s.email, 
			COALESCE(SUM(CASE WHEN a.status = 'Present' THEN 1 ELSE 0 END), 0) as present_count,
			COALESCE(SUM(CASE WHEN a.status = 'Absent' THEN 1 ELSE 0 END), 0) as absent_count
		FROM 
			students s
		LEFT JOIN 
			attendance a ON s.id = a.student_id AND a.date BETWEEN ? AND ?
		GROUP BY 
			s.id, s.name, s.email
		ORDER BY 
			s.name ASC
	`

	mock.ExpectQuery(regexp.QuoteMeta(expectedQuery)).
		WithArgs(startDate, endDate).
		WillReturnError(sql.ErrConnDone)

	reports, err := GetAttendanceReport(startDate, endDate)

	assert.Error(t, err)
	assert.Nil(t, reports)
	assert.NoError(t, mock.ExpectationsWereMet())
}
