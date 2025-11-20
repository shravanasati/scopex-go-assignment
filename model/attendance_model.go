package model

// Attendance struct
type Attendance struct {
	ID        int64     `json:"id" example:"1"`
	StudentID int64     `json:"student_id" example:"1" binding:"required"`
	Date      string    `json:"date" example:"2023-10-27" binding:"required"` // Using string for date input, could be time.Time
	Status    string    `json:"status" example:"Present" binding:"required,oneof=Present Absent"`
}

// Attendances array of Attendance type
type Attendances []Attendance
