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

	owm "github.com/briandowns/openweathermap"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	emj "github.com/kenshaw/emoji"
)

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
	EnableAutoMod  = flag.Bool("enable-automod", false, "Enable the auto-moderation")
)

var (
	//declaring secrets
	token        = ""
	aptly        = ""
	logChannelID = ""
	botPrefix    = ""
	buffer       = make([][]byte, 0)

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
			logCommand(s, i.GuildID, "time", i.Member.User)
		},
	}
)

var greetings = []string{
	"Hey",
	"It's good to see you again",
	"What's up?",
	"It's a pleasure to meet you",
}

var rpschoices = []string{
	"rock", "paper", "scissors",
}

func init() { flag.Parse() }

// Main function of the bot, called on startup.
func main() {
	//load env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
		return
	}
	//apply env variables to secrets
	token = os.Getenv("TOKEN")
	aptly = os.Getenv("API_KEY")
	logChannelID = os.Getenv("LOG_CHANNEL_ID")
	botPrefix = os.Getenv("BOT_PREFIX")

	// Load the sound file.
	err = loadSound("airhorn.dca")
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

	// Enable message tracking to the session.
	trackSessionMessages(dg)

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Setup intents
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers | discordgo.IntentsGuildPresences | discordgo.IntentsGuildVoiceStates | discordgo.IntentAutoModerationExecution | discordgo.IntentMessageContent

	if *EnableAutoMod {
		rule := initAutoModeration(dg)
		defer func() {
			fmt.Println("Removing AutoModerationRule")
			dg.AutoModerationRuleDelete(*GuildID, rule.ID)
		}()
	}

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

	// Register the guildMemberJoin func as a callback for GuildMemberAdd events.
	dg.AddHandler(guildMemberJoin)

	// Register the guildMemberLeave func as a callback for GuildMemberRemove events.
	dg.AddHandler(guildMemberLeave)

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

func logCommand(s *discordgo.Session, gID, command string, user *discordgo.User) {
	// Add slash prefix to the log if given command is a slash command
	// Else add prefix as BotPrefix
	var prefix = "Slash"
	if strings.Contains(command, botPrefix) {
		prefix = "BotPrefix"
	}

	// Log commands in the channel
	if logChannelID != "" {
		guild, err := s.Guild(gID)
		if err != nil {
			fmt.Printf("Error getting guild info: %s", err)
			return
		}

		s.ChannelMessageSend(
			logChannelID,
			fmt.Sprintf(`[%s] %s > %s > %s`, prefix, guild.Name, user, command),
		)
	}
}

// Sets the maximum number of messages to track for a session
func trackSessionMessages(dg *discordgo.Session) {
	dg.State.MaxMessageCount = 10000
}

func getGuild(s *discordgo.Session, channelID string) (*discordgo.Guild, error) {

	// Find the channel that the message came from.
	c, err := s.State.Channel(channelID)
	if err != nil {
		return nil, fmt.Errorf("Could not find channel")
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return nil, fmt.Errorf("Could not find guild")
	}
	return g, nil
}

// This handles all commands sent to the bot
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	// if the message starts with the prefix, then we know it's a command
	if strings.HasPrefix(m.Content, botPrefix) {
		switch strings.ToLower(strings.Split(m.Content, " ")[0]) {
		case botPrefix + "gollyhelp":
			s.ChannelMessageSend(m.ChannelID, formatHelpMessage())
		case botPrefix + "ping":
			s.ChannelMessageSend(m.ChannelID, "pong!")
		case botPrefix + "greet": // if the command is !greet
			s.ChannelMessageSend(m.ChannelID, randomGreeting(s, m))
		case botPrefix + "coinflip":
			s.ChannelMessageSend(m.ChannelID, coinFlip(s, m))
		case botPrefix + "horn":
			g, err := getGuild(s, m.ChannelID)
			if err != nil {
				fmt.Println("Error getting guild:", err)

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
		case botPrefix + "weather":
			w, err := owm.NewCurrent("F", "EN", aptly) // Returns weather in fahrenheit and English
			if err != nil {
				log.Fatalln(err)
			}
			var location = strings.Title(strings.ToLower(strings.SplitN(m.Content, " ", 2)[1]))
			w.CurrentByName(location)
			var result = fmt.Sprintf("Feels Like: %.2f°F\nTemperature: %.2f°F\nMin Temperature: %.2f°F\nMax Temperature: %.2f°F\nHumidity: %d%%\nWind speed: %.2fm/s\n", w.Main.FeelsLike, w.Main.Temp, w.Main.TempMin, w.Main.TempMax, w.Main.Humidity, w.Wind.Speed)
			for _, item := range w.Weather {
				result += fmt.Sprintf("%s: %s\n", item.Main, item.Description)
			}
			s.ChannelMessageSend(m.ChannelID, result)

		case botPrefix + "serverinfo":
			// sends embed containing server info
			s.ChannelMessageSendEmbed(m.ChannelID, serverinfo(s, m))
		case botPrefix + "remindme": //!remindme command
			var remindMessage = strings.SplitN(m.Content, " ", 3)[2]
			timer, err := strconv.Atoi(strings.SplitN(m.Content, " ", 3)[1])
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error! %s\nPlease use the correct syntax: !remindme <seconds> <message>", err))
			} else {
				s.ChannelMessageSend(m.ChannelID, "Reminder added!")
				s.ChannelMessageSend(m.ChannelID, remindMe(s, m, remindMessage, timer))
			}
		case botPrefix + "raffle":
			var msgContent = strings.Split(m.Content, " ")

			// Check if message ID and emoji are present
			if len(strings.Split(m.Content, " ")) < 3 {
				fmt.Println("Error finding message id or emoji")
				s.ChannelMessageSend(m.ChannelID, fmt.Sprint("Something went wrong, try again!\nPlease use the correct syntax : !raffle <message id> <reaction>\n"))
			}
			var messageID = msgContent[1]
			var emoji = msgContent[2]

			s.ChannelMessageSend(m.ChannelID, raffle(s, m.ChannelID, messageID, emoji))
		case botPrefix + "playrps":
			var rpsChoice = strings.ToLower(strings.Split(m.Content, " ")[1])
			if rps(s, m, rpsChoice) == "" {
				s.ChannelMessageSend(m.ChannelID, "Invalid choice! Use syntax `!playrps <choice>`.\nYour choices are: rock, paper and scissors")
			} else {
				s.ChannelMessageSend(m.ChannelID, rps(s, m, rpsChoice))
			}
		case botPrefix + "nickchange":
			st, err := s.UserChannelCreate(m.Author.ID)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("could not create a private channel for %v", m.Author))
			}
			nickArgs := strings.SplitN(m.Content, " ", 4)
			user := nickArgs[1]
			nick := nickArgs[2]
			g, err := getGuild(s, m.ChannelID)
			if err != nil {
				s.ChannelMessageSend(st.ID, fmt.Sprintf("an error getting guild: %v", err))
			}
			if err := s.GuildMemberNickname(g.ID, user, nick); err != nil {
				s.ChannelMessageSend(st.ID, fmt.Sprintf("an error occured while changing the nickname of %v to %v: %v", user, nick, err))
				return
			}
			s.ChannelMessageSend(st.ID, fmt.Sprintf("your nick has been changed"))
			return
		default:
			fmt.Println("Command not implemented")
			return
		}
		logCommand(s, m.GuildID, m.Content, m.Author)

		// if the message doesn't start with the prefix, then we check if it matches
		// one of the predefined messages to respond too
	} else {
		switch strings.Contains(strings.ToLower(m.Content), "hi golly") {
		case true:
			s.ChannelMessageSend(m.ChannelID, randomGreeting(s, m))
		}
	}
}

func formatHelpMessage() string {
	commands := map[string]string{
		"ping":                         "Reply with pong",
		"greet":                        "Reply with a greeting message",
		"coinflip":                     "Flip the coin!",
		"horn":                         "Honk the summoner",
		"weather <location>":           "How is the weather in <location> ? ",
		"remindme <seconds> <message>": "Create a reminder",
		"raffle <message id> <emoji>":  "Reply with a random user who reacted to given message with given emoji",
		"playrps <choice>":             "Play rock, paper and scissors",
	}

	helpMessage := "Available commands:\n"
	for cmd, desc := range commands {
		helpMessage += fmt.Sprintf("- `%s%s`: %s\n", botPrefix, cmd, desc)
	}

	return helpMessage
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

// !remindme command function
func remindMe(s *discordgo.Session, m *discordgo.MessageCreate, remindMessage string, timer int) string {
	<-time.After(time.Duration(timer) * time.Second)
	return fmt.Sprintf("%s! %s", m.Author.Mention(), "Reminder: "+remindMessage)
}

func initAutoModeration(session *discordgo.Session) *discordgo.AutoModerationRule {
	enabled := true
	rule, err := session.AutoModerationRuleCreate(*GuildID, &discordgo.AutoModerationRule{
		Name:        "GollyBot Auto Moderation",
		EventType:   discordgo.AutoModerationEventMessageSend,
		TriggerType: discordgo.AutoModerationEventTriggerKeyword,
		TriggerMetadata: &discordgo.AutoModerationTriggerMetadata{
			KeywordFilter: []string{"*Voldemort*"},
		},
		Actions: []discordgo.AutoModerationAction{
			{Type: discordgo.AutoModerationRuleActionBlockMessage},
		},

		Enabled: &enabled,
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully created the rule")

	session.AddHandlerOnce(func(s *discordgo.Session, e *discordgo.AutoModerationActionExecution) {
		_, err = session.AutoModerationRuleEdit(*GuildID, rule.ID, &discordgo.AutoModerationRule{
			Actions: []discordgo.AutoModerationAction{
				{Type: discordgo.AutoModerationRuleActionBlockMessage},
				{Type: discordgo.AutoModerationRuleActionSendAlertMessage, Metadata: &discordgo.AutoModerationActionMetadata{
					ChannelID: e.ChannelID,
				}},
			},
		})

		if err != nil {
			session.AutoModerationRuleDelete(*GuildID, rule.ID)
			panic(err)
		}

		s.ChannelMessageSend(e.ChannelID, "Shh we aren't supposed to speak that name!")
	})
	return rule
}

// raffle attempts to create list of users in the channel who reacted to given message
// with given emoji (or) reaction and returns a random user.
func raffle(s *discordgo.Session, channelID, messageID, emoji string) string {
	var reaction discordgo.Emoji

	// Find emoji id and name
	desc := emj.FromCode(emoji).Description
	emojis := emj.Gemoji()
	for _, em := range emojis {
		if em.Description == desc {
			reaction = discordgo.Emoji{
				ID:   em.Emoji,
				Name: em.Description,
			}
		}
	}

	// Find message with given id in the channel
	message, err := s.State.Message(channelID, messageID)
	if err != nil {
		//Could not find the message
		fmt.Println("Error finding message with ID :", messageID)
		return fmt.Sprint("Could not find message with ID : ", messageID)
	}

	// Find users who reacted to the message with given emoji
	reactedUsers, err := s.MessageReactions(channelID, messageID, reaction.ID, 100, "", "")
	if err != nil {
		fmt.Println(err)
		return fmt.Sprint("Could not find users reacted to this message")
	}

	// Returns, if no user reacted to the message
	if len(reactedUsers) == 0 {
		return fmt.Sprint("No user reacted to the message ", message.Content, "with emojiID ", reaction.ID)
	}

	// Pick random user from the list of users who reacted to the message
	rand.Seed(time.Now().UnixNano())
	winner := reactedUsers[rand.Intn(len(reactedUsers))]

	return fmt.Sprintf("The winner is: <@%s>", winner.ID)
}

func serverinfo(s *discordgo.Session, m *discordgo.MessageCreate) *discordgo.MessageEmbed {

	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("couldnt get the channel id")
	}

	// Find the guild for that channel.
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		log.Println("couldnt get the channel id")
	}

	owner := "<@" + g.OwnerID + ">"
	category_channel_count := 0
	text_channel_count := 0
	voice_channel_count := 0
	// finding the count of all channels
	for _, ch := range g.Channels {
		switch ch.Type {
		case discordgo.ChannelTypeGuildCategory:
			category_channel_count++
		case discordgo.ChannelTypeGuildText:
			text_channel_count++
		case discordgo.ChannelTypeGuildVoice:
			voice_channel_count++
		}
	}
	member_count := g.MemberCount
	role_count := len(g.Roles)
	embed := &discordgo.MessageEmbed{
		Title: "Server Info",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Server Owner",
				Value: owner,
			},
			{
				Name:  "Category Channels",
				Value: strconv.Itoa(category_channel_count),
			},
			{
				Name:  "Text Channels",
				Value: strconv.Itoa(text_channel_count),
			},
			{
				Name:  "Voice Channels",
				Value: strconv.Itoa(voice_channel_count),
			},
			{
				Name:  "Members",
				Value: strconv.Itoa(member_count),
			},
			{
				Name:  "Roles",
				Value: strconv.Itoa(role_count),
			},
		},
	}

	return embed
}

// !playrps command function
func rps(s *discordgo.Session, m *discordgo.MessageCreate, rpsChoice string) string {
	//pick a random rps choice
	randomIndex := rand.Intn(len(rpschoices))
	botChoice := rpschoices[randomIndex]
	//compare user choice with bot choice
	switch rpsChoice {
	case "rock":
		if botChoice == "rock" {
			return fmt.Sprintf("Your choice: %s\nMy choice: %s\n`It's a tie!`", rpsChoice, botChoice)
		}

		if botChoice == "paper" {
			return fmt.Sprintf("Your choice: %s\nMy choice: %s\nPaper covers rock. `I win!`", rpsChoice, botChoice)
		}

		if botChoice == "scissors" {
			return fmt.Sprintf("Your choice: %s\nMy choice: %s\nRock crushes scissors. `You win!`", rpsChoice, botChoice)
		}

	case "paper":
		if botChoice == "rock" {
			return fmt.Sprintf("Your choice: %s\nMy choice: %s\nPaper covers rock. `You win!`", rpsChoice, botChoice)
		}

		if botChoice == "paper" {
			return fmt.Sprintf("Your choice: %s\nMy choice: %s\n`It's a tie!`", rpsChoice, botChoice)

		}

		if botChoice == "scissors" {
			return fmt.Sprintf("Your choice: %s\nMy choice: %s\nScissors cut paper. `I win!`", rpsChoice, botChoice)
		}

	case "scissors":
		if botChoice == "rock" {
			return fmt.Sprintf("Your choice: %s\nMy choice: %s\nRock crushes scissors. `I win!`", rpsChoice, botChoice)
		}

		if botChoice == "paper" {
			return fmt.Sprintf("Your choice: %s\nMy choice: %s\nScissors cut paper. `You win!`", rpsChoice, botChoice)
		}

		if botChoice == "scissors" {
			return fmt.Sprintf("Your choice: %s\nMy choice: %s\n`It's a tie!`", rpsChoice, botChoice)
		}
	}
	return ""
}

// guildMemberJoin sends welcome message when a new member joins the server.
func guildMemberJoin(s *discordgo.Session, memberAdd *discordgo.GuildMemberAdd) {
	message := &discordgo.MessageSend{
		Content: fmt.Sprintf("Welcome to the channel %s", memberAdd.Mention()),
		Embeds: []*discordgo.MessageEmbed{{
			Image: &discordgo.MessageEmbedImage{
				URL:    "https://thumbs.dreamstime.com/b/welcome-banner-shiny-colorful-confetti-vector-paper-illustration-welcome-banner-colorful-confetti-100006906.jpg",
				Width:  10,
				Height: 10,
			},
		}},
	}

	//Find all channels in the guild.
	channels, _ := s.GuildChannels(memberAdd.GuildID)
	for _, c := range channels {
		// Send welcome message to general channel
		// Send message to a text channel, if there is no general channel in the guild
		if c.Name == "general" || c.Type == discordgo.ChannelTypeGuildText {
			s.ChannelMessageSendComplex(c.ID, message)
			return
		}
	}
	fmt.Sprint("There is no text channel to send the welcome message.")
	return
}

// guildMemberLeave sends goodbye message when a member leaves the server.
func guildMemberLeave(s *discordgo.Session, memberRemove *discordgo.GuildMemberRemove) {
	message := &discordgo.MessageSend{
		Content: fmt.Sprintf("Goodbye %s!", memberRemove.Mention()),
		Embeds: []*discordgo.MessageEmbed{{
			Image: &discordgo.MessageEmbedImage{
				URL:    "https://thumbs.dreamstime.com/b/goodbye-24589885.jpg",
				Width:  10,
				Height: 10,
			},
		}},
	}

	//Find all channels in the guild.
	channels, _ := s.GuildChannels(memberRemove.GuildID)
	for _, c := range channels {
		// Send goodbye message to general channel
		// Send message to a text channel, if there is no general channel in the guild
		if c.Name == "general" || c.Type == discordgo.ChannelTypeGuildText {
			s.ChannelMessageSendComplex(c.ID, message)
			return
		}
	}
	fmt.Sprint("There is no text channel to send the goodbye message.")
	return
}
