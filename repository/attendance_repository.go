package repository

import (
	"context"
	"log"
	"time"

	configuration "github.com/shravanasati/scopex-go-assignment/configuration"
	model "github.com/shravanasati/scopex-go-assignment/model"
)

// MarkAttendance records attendance for a student
func MarkAttendance(attendance model.Attendance) (int64, error) {
	db := configuration.DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := "INSERT INTO attendance (student_id, date, status) VALUES (?, ?, ?)"
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, attendance.StudentID, attendance.Date, attendance.Status)
	if err != nil {
		log.Println("Error marking attendance: " + err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetAttendanceByStudentID retrieves attendance records for a student
func GetAttendanceByStudentID(studentID int64) (model.Attendances, error) {
	db := configuration.DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var attendances model.Attendances

	query := "SELECT id, student_id, date, status FROM attendance WHERE student_id = ? ORDER BY date DESC"
	rows, err := db.QueryContext(ctx, query, studentID)
	if err != nil {
		log.Println("Error querying attendance: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a model.Attendance
		// Date comes as []uint8 or time.Time depending on driver.
		// Since we used ?parseTime=true, it should be time.Time, but we defined Date as string in struct.
		// We might need to scan into a temporary variable or change struct.
		// Let's scan into a string, MySQL driver usually handles this conversion if the target is string.
		err := rows.Scan(&a.ID, &a.StudentID, &a.Date, &a.Status)
		if err != nil {
			log.Println("Error scanning attendance: " + err.Error())
			return nil, err
		}
		attendances = append(attendances, a)
	}

	return attendances, nil
}

// GetAttendanceByDateRange retrieves attendance for a student within a date range
func GetAttendanceByDateRange(studentID int64, startDate, endDate string) (model.Attendances, error) {
	db := configuration.DB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var attendances model.Attendances

	query := "SELECT id, student_id, date, status FROM attendance WHERE student_id = ? AND date BETWEEN ? AND ? ORDER BY date ASC"
	rows, err := db.QueryContext(ctx, query, studentID, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var a model.Attendance
		err := rows.Scan(&a.ID, &a.StudentID, &a.Date, &a.Status)
		if err != nil {
			return nil, err
		}
		attendances = append(attendances, a)
	}

	return attendances, nil
}

// GetAttendanceReport retrieves aggregated attendance data for all students within a date range
func GetAttendanceReport(startDate, endDate string) (model.AttendanceReports, error) {
	db := configuration.DB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Longer timeout for report
	defer cancel()

	var reports model.AttendanceReports

	query := `
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

	rows, err := db.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		log.Println("Error querying attendance report: " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var r model.AttendanceReport
		err := rows.Scan(&r.StudentID, &r.StudentName, &r.StudentEmail, &r.PresentCount, &r.AbsentCount)
		if err != nil {
			log.Println("Error scanning attendance report: " + err.Error())
			return nil, err
		}
		reports = append(reports, r)
	}

	return reports, nil
}
