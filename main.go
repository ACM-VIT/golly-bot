package main

import (
	// Import the Discordgo package, and other required packages.

	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var (
	token     = "your token here"
	botPrefix = "!"
	commands  = []*discordgo.ApplicationCommand{
		{
			Name:        "time",
			Description: "return current time.",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"time": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Flags:   discordgo.MessageFlagsEphemeral,
					Content: "<t:" + strconv.FormatInt(time.Now().Unix(), 10) + ">",
				},
			})
		},
	}
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

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	fmt.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer dg.Close()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	if *RemoveCommands {
		fmt.Println("\nRemoving commands...")
		for _, v := range registeredCommands {
			err := dg.ApplicationCommandDelete(dg.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

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
