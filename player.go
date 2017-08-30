package main

import (
	"github.com/jonas747/dca"
	"log"
	"time"
	"io"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
)

type Player struct {
	Folder string
	session *discordgo.Session
}

var (
	stream *dca.StreamingSession
)

func (p *Player) Play(guildId, channelId string) {
	vc, err := p.session.ChannelVoiceJoin(guildId, channelId, false, true)
	if err != nil {
		fmt.Println("Error channel voice join:", err)
	}
	fmt.Println("Reading folder: ", p.Folder)
	files, _ := ioutil.ReadDir(p.Folder)
	for _, f := range files {
		fmt.Println("Play audio file: ", f.Name())
		playAudioFile(vc, fmt.Sprintf("%s/%s", p.Folder, f.Name()))
	}
	vc.Disconnect()
	return
}

func (p *Player) Pause() {
	//TODO: invalid request
	if (!stream.Paused()) {
		stream.SetPaused(true)
	}
}

func (p *Player) Resume() {
	//TODO: invalid request
	if (stream.Paused()) {
		stream.SetPaused(false)
	}
}

func playAudioFile(v *discordgo.VoiceConnection, filename string) {
	err := v.Speaking(true)
	if err != nil {
		log.Fatal("Failed setting speaking", err)
	}

	defer v.Speaking(false)

	opts := dca.StdEncodeOptions
	opts.RawOutput = true
	opts.Bitrate = 120

	encodeSession, err := dca.EncodeFile(filename, opts)
	if err != nil {
		log.Fatal("Failed creating an encoding session: ", err)
	}

	done := make(chan error)
	stream = dca.NewStream(encodeSession, v, done)

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case err := <-done:
			if err != nil && err != io.EOF {
				log.Fatal("An error occured: ", err)
			}
			encodeSession.Truncate()
			return
		case <-ticker.C:
			stats := encodeSession.Stats()
			playbackPosition := stream.PlaybackPosition()

			fmt.Printf("Playback: %10s, Transcode Stats: Time: %5s, Size: %5dkB, Bitrate: %6.2fkB, Speed: %5.1fx\r", playbackPosition, stats.Duration.String(), stats.Size, stats.Bitrate, stats.Speed)
		}
	}
}
