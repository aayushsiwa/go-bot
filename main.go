package main

import (
	"discord-bot/bot"
	"discord-bot/config"
	"discord-bot/services"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

func main() {
	cfg := config.Load()
	services.Cfg = cfg

	log.Println("ApplicationID:", cfg.ApplicationID)
	log.Println("GuildID:", cfg.GuildID)
	log.Println("Token length:", len(cfg.BotToken))

	dg, err := discordgo.New("Bot " + cfg.BotToken)

	dg.AddHandler(bot.HandleMessage)
	dg.AddHandler(bot.HandleInteraction) // ✅ important

	bot.Init()
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Bot is ready!")

		cmds, _ := s.ApplicationCommands(cfg.ApplicationID, cfg.GuildID)
		log.Println("Existing commands:", len(cmds))

		err := bot.RegisterSlashCommands(s, cfg.ApplicationID, cfg.GuildID)
		if err != nil {
			log.Printf("RegisterSlashCommands ERROR: %#v\n", err)
			log.Fatal(err)
		}
	})

	err = dg.Open()
	if err != nil {
		log.Printf("dg.Open ERROR: %#v\n", err)
		log.Fatal(err)
	}

	// httpapi.Discord = dg

	// // ✅ Start HTTP server (for interactions)
	// go func() {
	// 	http.HandleFunc("/interactions", httpapi.InteractionsHandler)
	// 	http.HandleFunc("/health", httpapi.Health)

	// 	log.Println("HTTP server running on :" + cfg.PORT)
	// 	if err := http.ListenAndServe(":"+cfg.PORT, nil); err != nil {
	// 		log.Fatalf("error %v", err)
	// 	}
	// }()

	// ✅ Worker system (unchanged)
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
