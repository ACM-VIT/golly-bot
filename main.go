package main

import (
	// Import the Discordgo package, and other required packages.

	"encoding/binary"
	"flag"
	"fmt"
	"io"
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
	token = "MTAyODI5OTU4ODc2MDc4NDkxNg.GSFhRJ.z-UNB6eexlMvMB0zFw_d1WY5eAkpJET9Aazr00"
	// token     = "your token here"
	botPrefix = "!"
	buffer    = make([][]byte, 0)

	commands = []*discordgo.ApplicationCommand{
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

	// Load the sound file.
	err := loadSound("airhorn.dca")
	if err != nil {
		fmt.Println("Error loading sound: ", err)
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Setup intents
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuildPresences | discordgo.IntentsGuildVoiceStates

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
		case "!horn":
			// Find the channel that the message came from.
			c, err := s.State.Channel(m.ChannelID)
			if err != nil {
				// Could not find channel.
				return
			}

			// Find the guild for that channel.
			g, err := s.State.Guild(c.GuildID)
			if err != nil {
				// Could not find guild.
				return
			}

			// Look for the message sender in that guild's current voice states.
			for _, vs := range g.VoiceStates {
				if vs.UserID == m.Author.ID {
					err = playSound(s, g.ID, vs.ChannelID)
					if err != nil {
						fmt.Println("Error playing sound:", err)
					}

					return
				}
			}
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
func coinFlip(s *discordgo.Session, m *discordgo.MessageCreate) string {
	coin := []string{
		"heads",
		"tails",
	}

	rand.Seed(time.Now().UnixNano())

	// flip the coin
	side := coin[rand.Intn(len(coin))]

	return fmt.Sprintf("Flipped the coin and you get : %s", side)
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string) (err error) {

	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}
	// Sleep for a specified amount of time before playing the sound
	time.Sleep(250 * time.Millisecond)
	// Start speaking.
	vc.Speaking(true)
	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}
	// Stop speaking
	vc.Speaking(false)
	// Sleep for a specificed amount of time before ending.
	time.Sleep(250 * time.Millisecond)
	// Disconnect from the provided voice channel.
	vc.Disconnect()
	return nil
}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}
	var opuslen int16
	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)
		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}
		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}
