package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var Token string

// Init function, called on startup, before the main function.
func init() {
	// Setup required variables and other objects here.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file.")
	}
}

// Main function of the bot, called on startup.
func main() {
	Token = os.Getenv("TOKEN")
	bot, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating discod session", err)
		return
	}

	bot.AddHandler(messageCreate)
	// We care about receiving message events.
	bot.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection", err)
		return
	}

	// Wait until CTRL-C or other term signal is received.
	fmt.Println("Bot is running. \nPress CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)
	<-sc

	// Close the websocket gracefully.
	bot.Close()
}

// MessageCreate is called whenever a new message is created on any channel that the authenticated bot has access to.
// Set parameters accordingly.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
