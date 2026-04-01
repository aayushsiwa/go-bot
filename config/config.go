package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT           string
	ApplicationID  string
	GuildID        string
	GuildChannelID string
	BotToken       string
	PublicKey      string
	ChannelID      string
	RSSFeedSize    int
}

var cfg *Config

func Load() *Config {
	// load ..env (ignore error in prod)
	_ = godotenv.Load()

	rssFeedSize, _ := strconv.Atoi(mustGetEnv("RSS_FEED_SIZE"))

	cfg = &Config{
		PORT:           os.Getenv("PORT"),
		ApplicationID:  mustGetEnv("APPLICATION_ID"),
		GuildID:        mustGetEnv("GUILD_ID"),
		GuildChannelID: mustGetEnv("GUILD_CHANNEL_ID"),
		BotToken:       mustGetEnv("BOT_TOKEN"),
		PublicKey:      mustGetEnv("PUBLIC_KEY"),
		ChannelID:      mustGetEnv("CHANNEL_ID"),
		RSSFeedSize:    rssFeedSize,
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
