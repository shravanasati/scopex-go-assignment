package cronjob

import (
	"log"

	"github.com/robfig/cron/v3"
	service "github.com/shravanasati/scopex-go-assignment/service"
)

// InitCron initializes and starts the cron scheduler
func InitCron() {
	c := cron.New()

	// Weekly Report - Every Sunday at midnight
	_, err := c.AddFunc("0 0 * * 0", func() {
		service.GenerateWeeklyReport()
	})
	if err != nil {
		log.Fatal("Error adding weekly cron job: ", err)
	}

	// Monthly Report - 1st of every month at midnight
	_, err = c.AddFunc("0 0 1 * *", func() {
		service.GenerateMonthlyReport()
	})
	if err != nil {
		log.Fatal("Error adding monthly cron job: ", err)
	}

	c.Start()
	log.Println("Cron scheduler started")
}
