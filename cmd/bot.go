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
	tele "gopkg.in/telebot.v3"
)

var (
	BotToken = os.Getenv("BOT_TOKEN")
)

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
			Token:  BotToken,
			Poller: &tele.LongPoller{Timeout: 10 * time.Second},
		}

		b, err := tele.NewBot(pref)

		if err != nil {
			log.Fatalf("Please check BOT_TOKEN env variable. %s", err)
			return
		}

		b.Handle("/hello", func(c tele.Context) error {
			return c.Send("Hello!")
		})

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
