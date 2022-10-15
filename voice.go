package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

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

func loadPresetDCAFiles() error {

	err := loadSound("airhorn.dca")
	if err != nil {
		fmt.Println("Error loading sound: ", err)
		return err
	}

	return nil
}
