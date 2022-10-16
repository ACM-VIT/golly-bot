package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

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
	fmt.Println("There is no text channel to send the welcome message.")
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
	fmt.Println("There is no text channel to send the goodbye message.")
}
