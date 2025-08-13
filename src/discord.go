package main

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	// commandPrefix = "!"
	GuildID = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
)

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "helloworld",
			Description: "Sends a hello world message",
		},
		{
			Name:        "ping",
			Description: "Ping the bot to check if it's online",
		},
		{
			Name:        "create-game",
			Description: "Create special channel for a game",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "game",
					Description: "The game to create a channel for",
					Required:    true,
				},
			},
		},
		{
			Name:        "set-timezone",
			Description: "Set your timezone for accurate time tracking",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "timezone",
					Description: "Your timezone (e.g., America/New_York, Europe/London, Asia/Tokyo)",
					Required:    true,
				},
			},
		},
	}
)

func discordAddHandlers(discord *discordgo.Session) {
	// discord.AddHandler(discordPrefixedCommands)

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"helloworld":   basicCommand,
		"create-game":  createGameCommand,
		"ping":         basicCommand,
		"set-timezone": setTimezoneCommand,
	}

	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func basicCommand(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("Basic Command executed")

	// Defer the response
	err := discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Println("Failed to defer interaction response:", err)
		return
	}

	_, err = discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "Hey there! Congratulations, you just executed your first slash command",
	})
	if err != nil {
		fmt.Println("Failed to send followup message:", err)
		return
	}
}

func createGameCommand(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("Command executed")

	// Defer the response
	err := discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Println("Failed to defer interaction response:", err)
		return
	}

	gameName := i.ApplicationCommandData().Options[0].StringValue()

	role, _ := createDiscordRole(gameName, discord, i.GuildID)

	channelID, err := createDiscordTextChannel(gameName, discord, i.GuildID, role.ID)
	if err != nil {
		fmt.Println("Failed to create channel:", err)
		return
	}

	fmt.Print("Channel created: ", channelID)
	// setDiscordPermissions(discord, channelID, "TheReds", DISCORD_ALLOW)
	// setDiscordPermissions(discord, channelID, DISCORD_ROLE, DISCORD_DENY)

	_, err = discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: "Creating a new game channel...",
	})
	if err != nil {
		fmt.Println("Failed to send followup message:", err)
		return
	}
}

// isValidTimezone checks if the given timezone string is valid
func isValidTimezone(tz string) bool {
	_, err := time.LoadLocation(tz)
	return err == nil
}

func setTimezoneCommand(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("Set timezone command executed")

	// Defer the response
	err := discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Println("Failed to defer interaction response:", err)
		return
	}

	// Get the timezone from the command options
	timezone := i.ApplicationCommandData().Options[0].StringValue()
	timezone = strings.TrimSpace(timezone)

	// Validate the timezone
	if !isValidTimezone(timezone) {
		_, err = discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Invalid timezone! Please use a valid timezone like:\n• America/New_York\n• Europe/London\n• Asia/Tokyo\n• UTC\n\nFor a full list, see: https://en.wikipedia.org/wiki/List_of_tz_database_time_zones",
		})
		if err != nil {
			fmt.Println("Failed to send error followup message:", err)
		}
		return
	}

	// Get the user's Discord ID
	discordID := i.Member.User.ID
	username := i.Member.User.Username
	discriminator := i.Member.User.Discriminator

	// Get or create the user first
	_, err = GetOrCreateUser(discordID, username, discriminator)
	if err != nil {
		fmt.Printf("Failed to get or create user: %v\n", err)
		_, followupErr := discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Failed to set timezone. Please try again later.",
		})
		if followupErr != nil {
			fmt.Println("Failed to send error followup message:", followupErr)
		}
		return
	}

	// Update the user's timezone
	err = UpdateUserTimezone(discordID, timezone)
	if err != nil {
		fmt.Printf("Failed to update user timezone: %v\n", err)
		_, followupErr := discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Failed to set timezone. Please try again later.",
		})
		if followupErr != nil {
			fmt.Println("Failed to send error followup message:", followupErr)
		}
		return
	}

	// Get current time in the user's timezone for confirmation
	loc, _ := time.LoadLocation(timezone)
	currentTime := time.Now().In(loc)

	_, err = discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("✅ Successfully set your timezone to **%s**!\n🕐 Your current time is: %s",
			timezone, currentTime.Format("Monday, January 2, 2006 at 3:04 PM MST")),
	})
	if err != nil {
		fmt.Println("Failed to send success followup message:", err)
		return
	}

	fmt.Printf("User %s (%s) set timezone to %s\n", username, discordID, timezone)
}

// func basicCommand(discord *discordgo.Session, i *discordgo.InteractionCreate) {
// 	fmt.Println("Command executed")
// 	discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
// 		Type: discordgo.InteractionResponseChannelMessageWithSource,
// 		Data: &discordgo.InteractionResponseData{
// 			Content: "Hey there! Congratulations, you just executed your first slash command",
// 		},
// 	})
// }

// func discordPrefixedCommands(discord *discordgo.Session, message *discordgo.MessageCreate) {
// 	if message.Content[:1] != commandPrefix || message.Content == "" {
// 		return
// 	}

// 	switch message.Content[1:] {
// 	case "helloworld":
// 		discord.ChannelMessageSend(message.ChannelID, "Hello, world!")
// 	default:
// 		discord.ChannelMessageSend(message.ChannelID, "Unknown command")
// 	}
// }
