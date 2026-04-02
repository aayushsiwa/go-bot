package httpapi

import (
	"bytes"
	"crypto/ed25519"
	"discord-bot/bot"
	"discord-bot/config"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var Discord *discordgo.Session

type Interaction struct {
	Type      int    `json:"type"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	Data      struct {
		Name    string `json:"name"`
		Options []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"options"`
	} `json:"data"`
}

func verify(signature, timestamp string, body []byte) bool {
	cfg := config.Get()
	sig, _ := hex.DecodeString(signature)
	pub, _ := hex.DecodeString(cfg.PublicKey)

	message := append([]byte(timestamp), body...)
	return ed25519.Verify(pub, message, sig)
}

func InteractionsHandler(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)

	// 🔐 VERIFY SIGNATURE
	signature := r.Header.Get("X-Signature-Ed25519")
	timestamp := r.Header.Get("X-Signature-Timestamp")

	if !verify(signature, timestamp, body) {
		http.Error(w, "invalid request signature", http.StatusUnauthorized)
		return
	}

	// restore body for JSON parsing
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	var interaction Interaction
	err := json.NewDecoder(r.Body).Decode(&interaction)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// ✅ PING RESPONSE
	if interaction.Type == 1 {
		_, err := w.Write([]byte(`{"type":1}`))
		if err != nil {
			return
		}
		return
	}

	var feed string
	var optsLog []string
	var args []string

	if interaction.Type == 2 {
		name := interaction.Data.Name

		for _, opt := range interaction.Data.Options {
			optsLog = append(optsLog, opt.Name+"="+opt.Value)
			if opt.Name == "feed" {
				feed = opt.Value
			}
			if feed == "" {
				feed = "tech" // default OR return error
			}

			args = append(args, feed)
		}

		log.Printf("Received command: %s | options: %s",
			name,
			strings.Join(optsLog, ", "),
		)

		ctx := &bot.Context{
			Session:   Discord,
			ChannelID: interaction.ChannelID,
			GuildID:   interaction.GuildID,
			Args:      args,
		}

		res, ok := bot.Execute(name, ctx)
		if !ok {
			res = "Unknown command"
		}

		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"type": 4,
			"data": map[string]string{
				"content": res,
			},
		})
		if err != nil {
			return
		}

		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
