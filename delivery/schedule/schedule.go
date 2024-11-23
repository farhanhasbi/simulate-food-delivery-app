package schedule

import (
	"food-delivery-apps/usecase"
	"log"

	"github.com/robfig/cron/v3"
)

func StartCronJob(userUc usecase.UserUseCase) {
	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("@every 10m", func() {
			deletedRows, err := userUc.CleanUpExpiredTokens()
			if err != nil {
					log.Printf("Error cleaning up tokens: %v\n", err.Error())
			} else if deletedRows > 0 {
					log.Printf("Expired tokens cleaned up: %d tokens deleted\n", deletedRows)
			} else {
					log.Println("No expired tokens found to clean up")
			}
	})

	if err != nil {
			log.Printf("Error scheduling cron job: %v\n", err.Error())
			return
	}

	c.Start()
	defer c.Stop()

	select{}
}

