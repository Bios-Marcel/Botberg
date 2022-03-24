package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	session := must(discordgo.New("Bot " + strings.TrimSpace(string(must(io.ReadAll(must(os.Open("token"))))))))
	session.Identify.Intents = discordgo.IntentsGuildMessages

	danielMinutesRegex := regexp.MustCompile("(\\d+(?:\\.\\d+)?)")
	session.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		//Don't wanna cause an endless loop of answering our own messages.
		if m.Author.ID == s.State.User.ID {
			return
		}

		lowercased := strings.ToLower(m.Content)
		if strings.Contains(lowercased, "daniel") {
			var danielMinutes time.Duration
			var input string
			if strings.Contains(lowercased, "minute") {
				minutes := danielMinutesRegex.FindString(lowercased)
				if len(minutes) == 0 {
					return
				}

				minutesParsed, errParse := strconv.ParseFloat(minutes, 64)
				if errParse != nil {
					return
				}

				input = fmt.Sprintf("%.2f daniel minute", minutesParsed)
				if minutesParsed >= 1 {
					input = input + "s"
				}
				danielMinutes = time.Duration(minutesParsed * 2.75 * float64(time.Minute))
			} else if strings.Contains(lowercased, "hour") {
				hours := danielMinutesRegex.FindString(lowercased)
				if len(hours) == 0 {
					return
				}

				hoursParsed, errParse := strconv.ParseFloat(hours, 64)
				if errParse != nil {
					return
				}

				input = fmt.Sprintf("%.2f daniel hour", hoursParsed)
				if hoursParsed > 1 {
					input = input + "s"
				}
				danielMinutes = time.Duration(hoursParsed * 2.75 * float64(time.Minute) * 60)
			} else {
				return
			}

			_, errSend := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s equals %s in normal human duration.", input, danielMinutes))
			if errSend != nil {
				log.Println("Error sending message:", errSend)
			}
		}
	})

	if errOpen := session.Open(); errOpen != nil {
		log.Fatalln("Error establishing websocket connection:", errOpen)
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}
