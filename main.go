package main

import (
	// Import the Discordgo package, and other required packages.
        "strings"
	"github.com/bwmarrin/discordgo"

)

//constants are variables that never change
const (
token = "token here"
botPrefix = "!"
)
// Main function of the bot, called on startup.

func main() {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}
	u, err := dg.User("@me")
	if err != nil {
		panic(err)
	}
	botID = u.ID

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuildPresences

	if err = dg.Open(); err != nil {
		panic(err)
	}

	fmt.Println("Bot Online! CTRL C to turn it off")
	select {}
}
// this handles all commands and makes it so when a player uses a command the bot actually runs it
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == botID {
		return
	}
	args := strings.Split(strings.TrimPrefix(m.Content, botPrefix), " ")
	command := args[0]
	if len(args) > 1 {
		args = args[1:]
	} else {
		args = nil
	}

	switch strings.ToLower(command) {
	

//Develop commands to the bot here.
     }
}
