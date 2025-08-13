package main

import (
	"flag"
	"fmt"
	"strconv"
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
		{
			Name:        "record-game",
			Description: "Record the result of a single game",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "leader",
					Description: "Your leader/character",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "opponent",
					Description: "Your opponent's leader/character",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: "Game category (Casual, Ranked, Locals, Tournament, etc.)",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Casual", Value: "Casual"},
						{Name: "Ranked", Value: "Ranked"},
						{Name: "Locals", Value: "Locals"},
						{Name: "Regional", Value: "Regional"},
						{Name: "National", Value: "National"},
						{Name: "Tournament", Value: "Tournament"},
						{Name: "Practice", Value: "Practice"},
						{Name: "Online", Value: "Online"},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "went_first",
					Description: "Did you go first?",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "won",
					Description: "Did you win?",
					Required:    true,
				},
			},
		},
		{
			Name:        "record-games",
			Description: "Record multiple games with the same leader",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "leader",
					Description: "Your leader/character for all games",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "category",
					Description: "Game category for all games (Casual, Ranked, Locals, Tournament, etc.)",
					Required:    false,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{Name: "Casual", Value: "Casual"},
						{Name: "Ranked", Value: "Ranked"},
						{Name: "Locals", Value: "Locals"},
						{Name: "Regional", Value: "Regional"},
						{Name: "National", Value: "National"},
						{Name: "Tournament", Value: "Tournament"},
						{Name: "Practice", Value: "Practice"},
						{Name: "Online", Value: "Online"},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "games",
					Description: "Games data: opponent1,first/second,win/loss;opponent2,first/second,win/loss",
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
		"record-game":  recordGameCommand,
		"record-games": recordGamesCommand,
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

func recordGameCommand(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("Record game command executed")

	// Defer the response
	err := discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Println("Failed to defer interaction response:", err)
		return
	}

	// Extract command options
	options := i.ApplicationCommandData().Options
	leader := ""
	opponent := ""
	category := "Casual" // Default category
	wentFirst := false
	won := false

	for _, option := range options {
		switch option.Name {
		case "leader":
			leader = option.StringValue()
		case "opponent":
			opponent = option.StringValue()
		case "category":
			category = NormalizeCategory(option.StringValue())
		case "went_first":
			wentFirst = option.BoolValue()
		case "won":
			won = option.BoolValue()
		}
	}

	// Get the user's Discord ID
	discordID := i.Member.User.ID
	username := i.Member.User.Username
	discriminator := i.Member.User.Discriminator

	// Get or create the user
	user, err := GetOrCreateUser(discordID, username, discriminator)
	if err != nil {
		fmt.Printf("Failed to get or create user: %v\n", err)
		_, followupErr := discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Failed to record game. Please try again later.",
		})
		if followupErr != nil {
			fmt.Println("Failed to send error followup message:", followupErr)
		}
		return
	}

	// Create the game result
	_, err = CreateGameResult(user.ID, leader, opponent, category, wentFirst, won)
	if err != nil {
		fmt.Printf("Failed to create game result: %v\n", err)
		_, followupErr := discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Failed to record game. Please try again later.",
		})
		if followupErr != nil {
			fmt.Println("Failed to send error followup message:", followupErr)
		}
		return
	}

	// Format response
	turnText := "second"
	if wentFirst {
		turnText = "first"
	}
	resultText := "lost"
	resultEmoji := "❌"
	if won {
		resultText = "won"
		resultEmoji = "✅"
	}

	_, err = discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: fmt.Sprintf("%s **Game Recorded!**\n🎮 **%s** vs **%s**\n📂 Category: **%s**\n🎯 Went **%s** • %s **%s**",
			resultEmoji, leader, opponent, category, turnText, resultEmoji, resultText),
	})
	if err != nil {
		fmt.Println("Failed to send success followup message:", err)
		return
	}

	fmt.Printf("User %s recorded game: %s vs %s (went %s, %s)\n", username, leader, opponent, turnText, resultText)
}

func recordGamesCommand(discord *discordgo.Session, i *discordgo.InteractionCreate) {
	fmt.Println("Record games command executed")

	// Defer the response
	err := discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
	if err != nil {
		fmt.Println("Failed to defer interaction response:", err)
		return
	}

	// Extract command options
	options := i.ApplicationCommandData().Options
	leader := ""
	category := "Casual" // Default category
	gamesData := ""

	for _, option := range options {
		switch option.Name {
		case "leader":
			leader = option.StringValue()
		case "category":
			category = NormalizeCategory(option.StringValue())
		case "games":
			gamesData = option.StringValue()
		}
	}

	// Get the user's Discord ID
	discordID := i.Member.User.ID
	username := i.Member.User.Username
	discriminator := i.Member.User.Discriminator

	// Get or create the user
	user, err := GetOrCreateUser(discordID, username, discriminator)
	if err != nil {
		fmt.Printf("Failed to get or create user: %v\n", err)
		_, followupErr := discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: "❌ Failed to record games. Please try again later.",
		})
		if followupErr != nil {
			fmt.Println("Failed to send error followup message:", followupErr)
		}
		return
	}

	// Parse games data
	// Expected format: opponent1,first/second,win/loss;opponent2,first/second,win/loss
	games := strings.Split(gamesData, ";")
	successCount := 0
	var gameResults []string

	for _, gameStr := range games {
		gameStr = strings.TrimSpace(gameStr)
		if gameStr == "" {
			continue
		}

		parts := strings.Split(gameStr, ",")
		if len(parts) != 3 {
			_, followupErr := discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: fmt.Sprintf("❌ Invalid game format: '%s'\nExpected format: opponent,first/second,win/loss", gameStr),
			})
			if followupErr != nil {
				fmt.Println("Failed to send error followup message:", followupErr)
			}
			return
		}

		opponent := strings.TrimSpace(parts[0])
		turnStr := strings.ToLower(strings.TrimSpace(parts[1]))
		resultStr := strings.ToLower(strings.TrimSpace(parts[2]))

		// Parse turn order
		var wentFirst bool
		if turnStr == "first" {
			wentFirst = true
		} else if turnStr == "second" {
			wentFirst = false
		} else {
			_, followupErr := discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: fmt.Sprintf("❌ Invalid turn format: '%s'\nUse 'first' or 'second'", turnStr),
			})
			if followupErr != nil {
				fmt.Println("Failed to send error followup message:", followupErr)
			}
			return
		}

		// Parse result
		var won bool
		if resultStr == "win" || resultStr == "won" {
			won = true
		} else if resultStr == "loss" || resultStr == "lost" || resultStr == "lose" {
			won = false
		} else {
			_, followupErr := discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: fmt.Sprintf("❌ Invalid result format: '%s'\nUse 'win/won' or 'loss/lost/lose'", resultStr),
			})
			if followupErr != nil {
				fmt.Println("Failed to send error followup message:", followupErr)
			}
			return
		}

		// Create the game result
		_, err = CreateGameResult(user.ID, leader, opponent, category, wentFirst, won)
		if err != nil {
			fmt.Printf("Failed to create game result: %v\n", err)
			_, followupErr := discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
				Content: fmt.Sprintf("❌ Failed to record game against %s. Please try again later.", opponent),
			})
			if followupErr != nil {
				fmt.Println("Failed to send error followup message:", followupErr)
			}
			return
		}

		successCount++

		// Format this game result
		turnText := "second"
		if wentFirst {
			turnText = "first"
		}
		resultText := "lost"
		resultEmoji := "❌"
		if won {
			resultText = "won"
			resultEmoji = "✅"
		}

		gameResults = append(gameResults, fmt.Sprintf("%s **%s** vs **%s** (went %s, %s)",
			resultEmoji, leader, opponent, turnText, resultText))
	}

	// Send success message
	responseContent := fmt.Sprintf("✅ **%s Games Recorded!**\n📂 Category: **%s**\n\n%s",
		strconv.Itoa(successCount), category, strings.Join(gameResults, "\n"))

	_, err = discord.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Content: responseContent,
	})
	if err != nil {
		fmt.Println("Failed to send success followup message:", err)
		return
	}

	fmt.Printf("User %s recorded %d games with leader %s\n", username, successCount, leader)
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
