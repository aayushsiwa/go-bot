package services

import (
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

func StartRSSCron(send func(string)) {
	loc, _ := time.LoadLocation("Asia/Kolkata")
	c := cron.New(cron.WithLocation(loc))

	lastSeen := make(map[string]string)

	c.AddFunc("0 5 * * *", func() {
		log.Println("Running RSS cron...")

		for name, url := range FeedMap {

			items, err := FetchLatest(url, 3)
			if err != nil || len(items) == 0 {
				continue
			}

			latest := items[0].Link

			// skip if already seen
			if lastSeen[name] == latest {
				continue // ✅ not return
			}

			lastSeen[name] = latest

			msg := fmt.Sprintf("📰 **%s (%s)**\n", name, url)

			for i, item := range items {
				msg += fmt.Sprintf("%d. [%s](%s)\n", i+1, item.Title, item.Link)
			}

			send(msg)
		}
	})

	c.Start()
}
