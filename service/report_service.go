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

	now := time.Now()
	endDate := now.Format("2006-01-02")
	startDate := now.AddDate(0, 0, -7).Format("2006-01-02")

	reports, err := repository.GetAttendanceReport(startDate, endDate)
	if err != nil {
		log.Println("Error fetching weekly attendance report: ", err)
		return
	}

	for _, report := range reports {
		output := fmt.Sprintf(
			"Weekly Report for %s (%s)\nPeriod: %s to %s\nPresent: %d, Absent: %d\n-----------------------------",
			report.StudentName, report.StudentEmail, startDate, endDate, report.PresentCount, report.AbsentCount,
		)
		fmt.Println(output)
	}
	log.Println("Weekly Attendance Report Generation Completed.")
}

// GenerateMonthlyReport generates and prints monthly attendance reports for all students
func GenerateMonthlyReport() {
	log.Println("Starting Monthly Attendance Report Generation...")

	now := time.Now()
	endDate := now.Format("2006-01-02")
	startDate := now.AddDate(0, -1, 0).Format("2006-01-02")

	reports, err := repository.GetAttendanceReport(startDate, endDate)
	if err != nil {
		log.Println("Error fetching monthly attendance report: ", err)
		return
	}

	for _, report := range reports {
		output := fmt.Sprintf(
			"Monthly Report for %s (%s)\nPeriod: %s to %s\nPresent: %d, Absent: %d\n-----------------------------",
			report.StudentName, report.StudentEmail, startDate, endDate, report.PresentCount, report.AbsentCount,
		)
		fmt.Println(output)
	}
	log.Println("Monthly Attendance Report Generation Completed.")
}
