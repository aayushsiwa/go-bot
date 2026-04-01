package bot

import (
	"discord-bot/services"
	"log"

	"github.com/bwmarrin/discordgo"
)

func Init() {

	Register(Command{
		Name:        "ping",
		Description: "Check if bot is alive",
		Execute: func(ctx *Context) string {
			return "Pong! 🏓"
		},
	})

	Register(Command{
		Name:        "sys",
		Description: "Get system stats",
		Execute: func(ctx *Context) string {
			return services.GetSystemStats()
		},
	})

	Register(Command{
		Name:        "rss",
		Description: "Get curated RSS feed",
		Execute: func(ctx *Context) string {
			log.Printf("Executing RSS command with args: %v", ctx.Args)
			if len(ctx.Args) == 0 {
				return "Please provide a feed (tech/gaming)"
			}

			feed := ctx.Args[0]

			go func() {
				feeds, err := services.GetFeedItems(feed)
				if err != nil {
					log.Println(err)
				}
				if err != nil || len(feeds) == 0 {
					return
				}

				SendRSS(ctx.Session, ctx.ChannelID, ctx.GuildID, feeds, feed)
			}()

			return "Sending you the latest " + feed + " news! 🚀"
		},
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "feed",
				Description: "Choose a feed",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{Name: "Tech", Value: "tech"},
					{Name: "News", Value: "news"},
					{Name: "Gaming", Value: "gaming"},
				},
			},
		},
	})
}

func RegisterSlashCommands(s *discordgo.Session, appID string, guildID string) {
	for _, cmd := range Registry {

		_, err := s.ApplicationCommandCreate(appID, guildID, &discordgo.ApplicationCommand{
			Name:        cmd.Name,
			Description: cmd.Description,
			Options:     cmd.Options,
		})
		if err != nil {
			return
		}

		log.Printf("Registered command: %s", cmd.Name)

	}
}
