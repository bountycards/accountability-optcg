package main

import (
	"flag"
	"fmt"

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
	}
)

func discordAddHandlers(discord *discordgo.Session) {
	// discord.AddHandler(discordPrefixedCommands)

	commandHandlers := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"helloworld":  basicCommand,
		"create-game": createGameCommand,
		"ping":        basicCommand,
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
