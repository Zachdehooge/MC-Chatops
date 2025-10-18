package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	h "github.com/zachdehooge/MC-Chatops/functions"
)

// Global Variables
var s *discordgo.Session

func init() {
	godotenv.Load()
	log.Print("Getting bot token from .env file")
	var BotToken = os.Getenv("TOKEN")
	var err error
	s, err = discordgo.New("Bot " + BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v | Check the .env", err)
	}
}

// Slash Commands
var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "botstatus",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "bot uptime",
		},
		{
			Name: "serverstatus",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "server uptime",
		},
		{
			Name:        "start",
			Description: "starts the minecraft server",
		},
		{
			Name:        "stop",
			Description: "stops the minecraft server",
		},
		{
			Name:        "scale",
			Description: "scales the minecraft server | default is auto",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"botstatus": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Bot Uptime: %s", h.Uptime()),
				},
			})
		},
		"serverstatus": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Server Uptime: %s\nServer Response Code: %v", h.Uptime(), h.ServerStatus()),
				},
			})
		},
		"start": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Starting the server",
				},
			})
		},
		"stop": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Stopping the server",
				},
			})
		},
		"scale": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Server scaling is set to: ",
				},
			})
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

func main() {
	var GuildID = os.Getenv("GuildID")

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// Clean up ALL old commands before re-registering
	existing, err := s.ApplicationCommands(s.State.User.ID, GuildID)
	if err != nil {
		log.Fatalf("Failed to list existing commands: %v", err)
	}

	for _, cmd := range existing {
		err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, cmd.ID)
		if err != nil {
			log.Printf("Failed to delete old command '%v': %v", cmd.Name, err)
		} else {
			log.Printf("Deleted old command: %v", cmd.Name)
		}
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	log.Println("Refreshing commands...")
	_, err = s.ApplicationCommandBulkOverwrite(s.State.User.ID, GuildID, commands)
	if err != nil {
		log.Fatalf("Cannot refresh commands: %v", err)
	}

	log.Println("Gracefully shutting down.")
}
