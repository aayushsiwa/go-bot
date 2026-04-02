package main

import (
	"discord-bot/bot"
	"discord-bot/config"
	"discord-bot/httpapi"
	"discord-bot/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

func main() {
	cfg := config.Load()

	dg, err := discordgo.New("Bot " + cfg.BotToken)
	if err != nil {
		log.Println(err)
	}

	dg.AddHandler(bot.HandleMessage)
	bot.Init()
	bot.RegisterSlashCommands(dg, cfg.GuildID, cfg.ApplicationID)

	err = dg.Open()
	if err != nil {
		log.Println(err)
	}

	httpapi.Discord = dg

	// ✅ Start HTTP server (for interactions)
	go func() {
		http.HandleFunc("/interactions", httpapi.InteractionsHandler)
		http.HandleFunc("/health", httpapi.Health)

		log.Println("HTTP server running on :" + cfg.PORT)
		if err := http.ListenAndServe(":"+cfg.PORT, nil); err != nil {
			log.Fatal(err)
		}
	}()

	// ✅ Worker system (unchanged)
	services.Cfg = cfg
	manager := services.NewManager()

	cpuWorker := services.NewCPUWorker(
		80,
		30*time.Second,
		2*time.Minute,
		func(msg string) {
			_, err := dg.ChannelMessageSend(cfg.GuildChannelID, msg)
			if err != nil {
				return
			}
		},
	)

	services.StartRSSCron(func(feed string, items []*gofeed.Item) {
		bot.SendRSS(dg, cfg.GuildChannelID, cfg.GuildID, items, feed)
	})

	manager.Add(cpuWorker)
	manager.StartAll()

	log.Println("Bot + HTTP server running...")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	err = dg.Close()
	if err != nil {
		return
	}
	manager.StopAll()
}
