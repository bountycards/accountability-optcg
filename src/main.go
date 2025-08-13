package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func main() {
	discord, err := discordgo.New("Bot " + getEnv("DISCORD_TOKEN"))
	if err != nil {
		err_msg := "Error creating Discord session: " + err.Error()
		panic(err_msg)
	}

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	discordAddHandlers(discord)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		fmt.Printf("Registering command '%v'...\n", v.Name)
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	// panic(1)
	fmt.Println("Bot is now running. Press Ctrl+C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	defer removeCommands(registeredCommands, discord)

	defer discord.Close()

}

func removeCommands(registeredCommands []*discordgo.ApplicationCommand, discord *discordgo.Session) {
	fmt.Println("Removing commands...")
	for _, v := range registeredCommands {
		err := discord.ApplicationCommandDelete(discord.State.User.ID, *GuildID, v.ID)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
