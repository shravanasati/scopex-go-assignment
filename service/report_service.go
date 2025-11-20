package service

import (
	"fmt"
	"log"
	"time"

	repository "github.com/shravanasati/scopex-go-assignment/repository"
)

// GenerateWeeklyReport generates and prints weekly attendance reports for all students
func GenerateWeeklyReport() {
	log.Println("Starting Weekly Attendance Report Generation...")

	students, err := repository.StudentRepo.GetAllStudents(1000, 0) // Assuming < 1000 students for now, or implement pagination loop
	if err != nil {
		log.Println("Error fetching students for report: ", err)
		return
	}

	now := time.Now()
	endDate := now.Format("2006-01-02")
	startDate := now.AddDate(0, 0, -7).Format("2006-01-02")

	for _, student := range students {
		attendances, err := repository.GetAttendanceByDateRange(student.ID, startDate, endDate)
		if err != nil {
			log.Printf("Error fetching attendance for student %d: %v\n", student.ID, err)
			continue
		}

		presentCount := 0
		absentCount := 0
		for _, a := range attendances {
			if a.Status == "Present" {
				presentCount++
			} else {
				absentCount++
			}
		}

		report := fmt.Sprintf(
			"Weekly Report for %s (%s)\nPeriod: %s to %s\nPresent: %d, Absent: %d\n-----------------------------",
			student.Name, student.Email, startDate, endDate, presentCount, absentCount,
		)
		fmt.Println(report)
	}
	log.Println("Weekly Attendance Report Generation Completed.")
}

// GenerateMonthlyReport generates and prints monthly attendance reports for all students
func GenerateMonthlyReport() {
	log.Println("Starting Monthly Attendance Report Generation...")

	students, err := repository.StudentRepo.GetAllStudents(1000, 0)
	if err != nil {
		log.Println("Error fetching students for report: ", err)
		return
	}

	now := time.Now()
	endDate := now.Format("2006-01-02")
	startDate := now.AddDate(0, -1, 0).Format("2006-01-02")

	for _, student := range students {
		attendances, err := repository.GetAttendanceByDateRange(student.ID, startDate, endDate)
		if err != nil {
			log.Printf("Error fetching attendance for student %d: %v\n", student.ID, err)
			continue
		}

		presentCount := 0
		absentCount := 0
		for _, a := range attendances {
			if a.Status == "Present" {
				presentCount++
			} else {
				absentCount++
			}
		}

		report := fmt.Sprintf(
			"Monthly Report for %s (%s)\nPeriod: %s to %s\nPresent: %d, Absent: %d\n-----------------------------",
			student.Name, student.Email, startDate, endDate, presentCount, absentCount,
		)
		fmt.Println(report)
	}
	log.Println("Monthly Attendance Report Generation Completed.")
}
