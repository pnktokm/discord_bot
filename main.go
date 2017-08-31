package main

import (
	"fmt"
	"time"
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"math/rand"
	"os"
	"syscall"
	"os/signal"
)

const (
	Token = "Bot <TO_BE_DEFINED>"
	Folder = "<path>"
	fqdn    = "https://danbooru.donmai.us"
	popular = fqdn + "/explore/posts/popular"
)

var (
	player *Player
	httpClient *http.Client
)

func main() {
	discord, err := discordgo.New(Token)
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
	}
	httpClient = http.DefaultClient

	discord.AddHandler(onMessageCreate)
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Listening...")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discord.Close()
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println("Error getting channel: ", err)
		return
	}
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		return
	}
	fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

	switch {
	case strings.HasPrefix(m.Content, fmt.Sprintf("%s", "!help")):
		sender := &Sender{httpClient, s}
		commands := map[string]string{
			"!image": "適当な画像 3件を表示します。",
			"!music": "プレイリストを再生します。",
			"!pause": "再生中の曲を一時停止します。",
			"!resume": "一時停止中の曲を再生します。",
			"!kill": "botを停止します。"}
		var sendText = "【使い方】"
		for command, desc := range commands {
			sendText += fmt.Sprintf("\n%s: %s", command, desc)
		}
		sender.SendMessage(c.ID, sendText)

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s", "!image")):
		sender := &Sender{httpClient, s}
		doc, err := goquery.NewDocument(popular)
		if err != nil {
			panic(err)
		}
		now := time.Now().UTC().Format("2006-01-02")
		var arr []string
		doc.Find("#a-popular article").Each(func(_ int, s *goquery.Selection) {
			a, _ := s.Attr("data-file-url")
			arr = append(arr, fqdn + a)
		})
		doc2, _ := goquery.NewDocument(popular + "?date=" + now + "&scale=week")
		doc2.Find("#a-popular article").Each(func(_ int, s *goquery.Selection) {
			a, _ := s.Attr("data-file-url")
			arr = append(arr, fqdn + a)
		})
		siko := []string{ // !?
			"http://s1.dmcdn.net/fakq/1280x720-Krq.jpg",
			"https://i.ytimg.com/vi/CbCJK-93ubI/maxresdefault.jpg",
			"https://i.ytimg.com/vi/80tIjMrd7_c/hqdefault.jpg",
			"http://s1.dmcdn.net/Ar-r7/1280x720-9iB.jpg"}
		for _, s := range siko {
			arr = append(arr, s)
		}
		rand.Seed(time.Now().UnixNano())
		for i, _ := range arr {
			if (i < 3) {
				sender.SendFile(c.ID, arr[rand.Intn(len(arr)-1)])
				time.Sleep(500 * time.Millisecond)
			} else {
				return
			}
		}

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s", "!music")):
		player = &Player{Folder, s}

		for _, vs := range g.VoiceStates {
			if vs.UserID == m.Author.ID {
				player.Play(g.ID, vs.ChannelID)
			}
		}

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s", "!pause")):
		player.Pause()

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s", "!resume")):
		player.Resume()

	case strings.HasPrefix(m.Content, fmt.Sprintf("%s", "!kill")):
		s.Close()
		os.Exit(1)

	}
}
