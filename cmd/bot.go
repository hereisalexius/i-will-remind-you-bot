/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"gopkg.in/robfig/cron.v2"
	tele "gopkg.in/telebot.v3"
)

var (
	TeleToken     = os.Getenv("TELE_TOKEN")
	messagesCache = make(map[string]*MessageToRemind)
)

type MessageToRemind struct {
	//I'm lazy to make a proper handler, lets store whole context to send remind messages
	UserContext tele.Context
	MessageText string
	RimindTime  time.Time
}

// botCmd represents the bot command
var botCmd = &cobra.Command{
	Use:     "bot",
	Aliases: []string{"start"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		pref := tele.Settings{
			URL:    "",
			Token:  TeleToken,
			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		}

		b, err := tele.NewBot(pref)

		if err != nil {
			log.Fatalf("Please check BOT_TOKEN env variable. %s", err)
			return
		}

		b.Handle("/set", func(c tele.Context) error {
			var senderUsername = c.Sender().Username

			messagesCache[senderUsername] = new(MessageToRemind)
			messagesCache[senderUsername].UserContext = c

			log.Printf("New reminder has been initiated for %s\n", senderUsername)

			return c.Send("What I should remind you?")
		})

		b.Handle("/dismiss", func(c tele.Context) error {
			var senderUsername = c.Sender().Username

			messagesCache[senderUsername] = nil

			log.Printf("Reminder for %s dismissed.\n", senderUsername)

			return c.Send("Reminder dismissed. Would you like to /set new?")
		})

		b.Handle("/ping", func(c tele.Context) error {
			var senderUsername = c.Sender().Username

			if messagesCache[senderUsername] == nil || messagesCache[senderUsername].RimindTime.IsZero() {
				return c.Send("Reminder not set.")
			}

			return c.Send(fmt.Sprintf("Current reminder: %s\n"+
				"Time: %s\n"+
				"Would you like to /dismiss it?",
				messagesCache[senderUsername].MessageText,
				messagesCache[senderUsername].RimindTime.Format(time.Kitchen)))
		})

		b.Handle(tele.OnText, func(c tele.Context) error {
			var senderUsername = c.Sender().Username
			payload := c.Message().Payload
			text := c.Text()

			log.Println(payload, c.Text())

			var remindMsg = messagesCache[senderUsername]

			if remindMsg == nil {
				log.Printf("%s has not set any reminders. Treated as note.\n", senderUsername)
				return err
			}

			if remindMsg.MessageText == "" {
				remindMsg.MessageText = text
				log.Printf("%s has set message text: %s\n", senderUsername, remindMsg.MessageText)
				return c.Send("After how long to warn you?\n*Please set time duration like : 1h2m30s")
			}

			if remindMsg.RimindTime.IsZero() {
				reminderDuration, err := time.ParseDuration(text)

				if err != nil {
					return c.Send("Unable to set duration. Please try another one.")
				}
				remindMsg.RimindTime = time.Now().Add(reminderDuration)

				log.Printf("%s has set time to: %s\n", senderUsername, remindMsg.RimindTime)
				log.Printf("Reminder for %s created!\n", senderUsername)
				return c.Send("Done. You can /ping for status.")
			}

			return err
		})

		runCronJob()
		fmt.Printf("i-will-remind-you-bot %s started\n", version)
		b.Start()
	},
}

func init() {
	rootCmd.AddCommand(botCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// botCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// botCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func runCronJob() {
	job := cron.New()

	job.AddFunc("@every 5s", func() {
		for _, element := range messagesCache {
			if element != nil && !element.RimindTime.IsZero() {
				if time.Now().After(element.RimindTime) {
					element.UserContext.Send(fmt.Sprintf("Notifying you about %s!\n/dismiss ?", element.MessageText))
				}
			}
		}
	})

	job.Start()

}
