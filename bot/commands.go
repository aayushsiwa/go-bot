package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/mmcdole/gofeed"
)

type Context struct {
	Session   *discordgo.Session
	ChannelID string
	GuildID   string
	Args      []string
}

type Command struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	Execute     func(ctx *Context) string
}
type CommandFunc func(args []string) string

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if !strings.HasPrefix(m.Content, "!") {
		return
	}

	parts := strings.Split(m.Content[1:], " ")
	name := parts[0]
	args := parts[1:]

	log.Printf("[%s] %s", m.Author.Username, name)

	ctx := &Context{
		Session:   s,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
		Args:      args,
	}

	if res, ok := Execute(name, ctx); ok {
		_, err := s.ChannelMessageSend(m.ChannelID, res)
		if err != nil {
			return
		}
	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, "❌ Unknown command")
		if err != nil {
			return
		}
	}
}

func SendRSSAsMessages(s *discordgo.Session, channelID string, items []*gofeed.Item) {
	for i, item := range items {
		content := fmt.Sprintf("%d. [%s](%s)", i+1, item.Title, item.Link)
		_, err := s.ChannelMessageSend(channelID, content)
		if err != nil {
			return
		}
	}
}

func SendRSSInThread(s *discordgo.Session, channelID string, items []*gofeed.Item, feedName string) {

	// 1. Send base message
	msg, err := s.ChannelMessageSend(channelID, strings.ToUpper(feedName))
	if err != nil {
		return
	}

	// 2. Create thread
	thread, err := s.MessageThreadStart(channelID, msg.ID, "Feed -> ", 60)
	if err != nil {
		return
	}

	// 3. Send items in thread
	for _, item := range items {
		if _, err := s.ChannelMessageSendEmbed(thread.ID, &discordgo.MessageEmbed{
			Title: item.Title,
			URL:   item.Link,
		}); err != nil {
			return
		}
	}
}

func SendRSS(s *discordgo.Session, channelID, guildID string, items []*gofeed.Item, feed string) {

	if len(items) == 0 {
		_, err := s.ChannelMessageSend(channelID, "❌ No items found")
		if err != nil {
			return
		}
		return
	}

	if guildID == "" {
		SendRSSAsMessages(s, channelID, items)
		return
	}

	SendRSSInThread(s, channelID, items, feed)
}

func HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	// Convert options → args
	var args []string
	for _, opt := range data.Options {
		args = append(args, opt.StringValue())
	}

	ctx := &Context{
		Session:   s,
		ChannelID: i.ChannelID,
		GuildID:   i.GuildID,
		Args:      args,
	}

	// 🔥 Reuse your existing system
	if res, ok := Execute(data.Name, ctx); ok {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: res,
			},
		})
		if err != nil {
			log.Println("interaction error:", err)
		}
	} else {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Unknown command",
			},
		})
	}
}
