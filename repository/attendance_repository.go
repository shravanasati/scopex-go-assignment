package repository

import (
	"log"

	configuration "github.com/shravanasati/scopex-go-assignment/configuration"
	model "github.com/shravanasati/scopex-go-assignment/model"
)

// MarkAttendance records attendance for a student
func MarkAttendance(attendance model.Attendance) (int64, error) {
	db := configuration.DB
	query := "INSERT INTO attendance (student_id, date, status) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(attendance.StudentID, attendance.Date, attendance.Status)
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
	var attendances model.Attendances

	query := "SELECT id, student_id, date, status FROM attendance WHERE student_id = ? ORDER BY date DESC"
	rows, err := db.Query(query, studentID)
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
	var attendances model.Attendances

	query := "SELECT id, student_id, date, status FROM attendance WHERE student_id = ? AND date BETWEEN ? AND ? ORDER BY date ASC"
	rows, err := db.Query(query, studentID, startDate, endDate)
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
