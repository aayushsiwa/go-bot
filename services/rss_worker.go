package services

import (
	"discord-bot/config"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"
)

var Cfg *config.Config

func StartRSSCron(send func(feed string, items []*gofeed.Item)) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		log.Fatal(err)
	}
	c := cron.New(cron.WithLocation(loc))

	lastSeen := make(map[string]string)

	_, err = c.AddFunc(Cfg.RSSCronSchedule, func() {
		log.Println("Running RSS cron...")

		for name, url := range FeedMap {

			items, err := FetchLatest(url, Cfg.RSSFeedSize)
			if err != nil || len(items) == 0 {
				continue
			}

			latest := items[0].Link

			if lastSeen[name] == latest {
				continue
			}

			lastSeen[name] = latest

			// ✅ call injected function
			send(name, items)
		}
	})
	if err != nil {
		return
	}

	c.Start()
}
