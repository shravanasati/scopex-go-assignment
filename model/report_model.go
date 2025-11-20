package model

// AttendanceReport struct to hold aggregated report data
type AttendanceReport struct {
	StudentID    int64  `json:"student_id"`
	StudentName  string `json:"student_name"`
	StudentEmail string `json:"student_email"`
	PresentCount int    `json:"present_count"`
	AbsentCount  int    `json:"absent_count"`
}

// AttendanceReports array of AttendanceReport
type AttendanceReports []AttendanceReport
