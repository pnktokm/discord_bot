package main

import (
	"net/http"
	"github.com/bwmarrin/discordgo"
	"log"
	"fmt"
)

type Sender struct {
	HttpClient *http.Client
	session *discordgo.Session
}

func (s *Sender) SendMessage(channelId, msg string) {
	_, err := s.session.ChannelMessageSend(channelId, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func (s *Sender) SendFile(channelId, msg string) {
	resp, httpErr := s.HttpClient.Get(msg)
	defer resp.Body.Close()
	if httpErr != nil {
		fmt.Println(fmt.Sprintf("%s: %s", "http error", httpErr.Error()))
	}

	_, err := s.session.ChannelFileSend(channelId, msg, resp.Body)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}
