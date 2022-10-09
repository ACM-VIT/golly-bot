package main

import (
	// Import the Discordgo package, and other required packages.

	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

//constants are variables that never change
const (
	token     = "token here"
	botPrefix = "!"
)

var greetings = []string{
	"Hey",
	"It's good to see you again",
	"What's up?",
	"It's a pleasure to meet you",
}

// Main function of the bot, called on startup.
func main() {
	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Setup intents
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuildPresences

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This handles all commands sent to the bot
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// if the message starts with the prefix, then we know it's a command
	if strings.HasPrefix(m.Content, botPrefix) {
		switch strings.ToLower(strings.Split(m.Content, " ")[0]) {
		// if the command is !greet
		case "!greet":
			s.ChannelMessageSend(m.ChannelID, randomGreeting(s, m))
		case "!coinflip":
			s.ChannelMessageSend(m.ChannelID, coinFlip(s, m))
		}
		// if the message doesn't start with the prefix, then we check if it matches
		// one of the predefined messages to respond too
	} else {
		switch strings.Contains(strings.ToLower(m.Content), "hi golly") {
		case true:
			s.ChannelMessageSend(m.ChannelID, randomGreeting(s, m))
		}
	}
}

func randomGreeting(s *discordgo.Session, m *discordgo.MessageCreate) (greeting string) {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())
	// Shuffle the greetings to return a random greeting
	rand.Shuffle(len(greetings), func(i, j int) {
		greetings[i], greetings[j] = greetings[j], greetings[i]
	})

	return fmt.Sprintf("%s, %s", greetings[0], m.Author.Mention())
}

// Random coinflip command
func coinFlip(s *discordgo.Session, m *discordgo.MessageCreate) (string) {
	coin := []string{
                 "heads",
                 "tails",
         }

         rand.Seed(time.Now().UnixNano())

         // flip the coin
         side := coin[rand.Intn(len(coin))]

         return fmt.Sprintf("Flipped the coin and you get : %s", side)
}
