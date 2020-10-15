package init

import (
	"os"
	"time"

	"github.com/Aiscom-LLC/meals-api/backups"

	"github.com/Aiscom-LLC/meals-api/config"
	"github.com/Aiscom-LLC/meals-api/repository"
	"github.com/Aiscom-LLC/meals-api/utils"
)

func init() {
	orderRepo := repository.NewOrderRepo()
	clientRepo := repository.NewClientRepo()
	clients, _ := clientRepo.GetAll()

	for _, client := range clients {
		if client.AutoApproveOrders {
			_, _ = clientRepo.InitAutoApprove(client.ID.String())
		}
	}

	config.CRON.Cron.Start()
	_ = config.CRON.Cron.AddFunc("@every 0h1m0s", func() {
		currentDay := utils.GetCurrentDay()
		currentTime := time.Now().Format("15:04")
		for _, entry := range config.CRON.Entries {
			for entryKey, clientIDMap := range entry {
				for key, value := range clientIDMap {
					if currentDay == key && currentTime == value {
						nextDay := time.Now().Add(time.Hour * 24).UTC().Truncate(time.Hour * 24).Format(time.RFC3339)
						_ = orderRepo.ApproveOrders(entryKey, nextDay)
					}
				}
			}
		}
	})
	if os.Getenv("BACKUP") == "true" {
		_ = config.CRON.Cron.AddFunc(utils.CronStringCreator("Europe/Moscow", "00", "00"), backups.CreateBackup)
	}
}
