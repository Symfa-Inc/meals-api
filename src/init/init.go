package init

import (
	"github.com/Aiscom-LLC/meals-api/src/config"
	"github.com/Aiscom-LLC/meals-api/src/repository"
	"github.com/Aiscom-LLC/meals-api/src/utils"
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
			for entryKey, clientIDMap := range entry {
				for key, value := range clientIDMap {
					if currentDay == key && currentTime == value {
						nextDay := time.Now().Add(time.Hour * 24).UTC().Truncate(time.Hour * 24).Format(time.RFC3339)
						orderRepo.ApproveOrders(entryKey, nextDay)
					}
				}
			}
		}
	})
}
