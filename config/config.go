package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	PORT            string
	ApplicationID   string
	GuildID         string
	GuildChannelID  string
	BotToken        string
	PublicKey       string
	ChannelID       string
	RSSFeedSize     int
	RSSCronSchedule string
}

var cfg *Config

func Load() *Config {
	// load ..env (ignore error in prod)
	// _ = godotenv.Load()

	rssFeedSize, err := strconv.Atoi(mustGetEnv("RSS_FEED_SIZE"))
	if err != nil {
		log.Fatal(err)
	}

	schedule := os.Getenv("RSS_CRON_SCHEDULE")
	if schedule == "" {
		schedule = "0 5 * * *" // fallback
	}

	cfg = &Config{
		PORT:            mustGetEnv("PORT"),
		ApplicationID:   mustGetEnv("APPLICATION_ID"),
		GuildID:         mustGetEnv("GUILD_ID"),
		GuildChannelID:  mustGetEnv("GUILD_CHANNEL_ID"),
		BotToken:        mustGetEnv("BOT_TOKEN"),
		PublicKey:       mustGetEnv("PUBLIC_KEY"),
		ChannelID:       mustGetEnv("CHANNEL_ID"),
		RSSFeedSize:     rssFeedSize,
		RSSCronSchedule: schedule,
	}

	return cfg
}

func Get() *Config {
	return cfg
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Missing required .env variable: %s", key)
	}
	return val
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
