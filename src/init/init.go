package init

import (
	"go_api/src/config"
	"go_api/src/repository"
	"go_api/src/utils"
	"time"
)

func init() {
	orderRepo := repository.NewOrderRepo()
	clientRepo := repository.NewClientRepo()
	clients, _ := clientRepo.GetAll()

	for _, client := range clients {
		if client.AutoApproveOrders {
			clientRepo.InitAutoApprove(client.ID.String())
		}
	}

	config.CRON.Cron.Start()
	config.CRON.Cron.AddFunc("@every 0h1m0s", func() {
		currentDay := utils.GetCurrentDay()
		currentTime := time.Now().Format("15:04")
		for _, entry := range config.CRON.Entries {
			for entryKey, clientIdMap := range entry {
				for key, value := range clientIdMap {
					if currentDay == key && currentTime == value {
						nextDay := time.Now().Add(time.Hour * 24).UTC().Truncate(time.Hour * 24).Format(time.RFC3339)
						orderRepo.ApproveOrders(entryKey, nextDay)
					}
				}
			}
		}
	})
}
