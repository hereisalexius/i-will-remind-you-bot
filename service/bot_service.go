package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/robfig/cron.v2"
	telebot "gopkg.in/telebot.v3"
)

type MessageToRemind struct {
	UserContext telebot.Context
	MessageText string
	RimindTime  time.Time
}

var (
	TeleToken      = os.Getenv("TELE_TOKEN")
	messagesCache  = make(map[string]*MessageToRemind)
	descriptionMsg = "I can /help you to /set simple reminder" +
		" and will annoy you with notifications so you wont forget the thing!\n" +
		"\n/start - obviously, to start using bot" +
		"\n/set - you will be prompted to questionnaire." +
		"\n		Question #1: What I should remind you? - Just type what you want to be reminded of. (example: Turn off oven)" +
		"\n		Question #2: After how long to warn you? - Provide duration for time when you should be notified (example: 30m)" +
		"\n/dismiss - to dismiss notification on any stage" +
		"\n/ping - to check if notification was set" +
		"\n/help - to see help"
)

func StartBot(appVersion string) {
	pref := telebot.Settings{
		URL:    "",
		Token:  TeleToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := telebot.NewBot(pref)

	log.Println("Bot instance created.")

	if err != nil {
		log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
		return
	}

	initStartHandler(*b)
	initHelpHandler(*b)
	initSetHandler(*b)
	initDismissHandler(*b)
	initPingHandler(*b)
	initOnTextHandler(*b, err)

	runReminderCronJob()
	fmt.Printf("i-will-remind-you-bot %s started.\n", appVersion)
	b.Start()
}

func initSetHandler(b telebot.Bot) {
	b.Handle("/set", func(c telebot.Context) error {
		var senderUsername = c.Sender().Username

		messagesCache[senderUsername] = new(MessageToRemind)
		messagesCache[senderUsername].UserContext = c

		log.Printf("New reminder has been initiated for %s\n", senderUsername)

		return c.Send("What I should remind you?")
	})
	log.Println("Handler for /set - initialized.")
}

func initDismissHandler(b telebot.Bot) {
	b.Handle("/dismiss", func(c telebot.Context) error {
		var senderUsername = c.Sender().Username

		messagesCache[senderUsername] = nil

		log.Printf("Reminder for %s dismissed.\n", senderUsername)

		return c.Send("Reminder dismissed. Would you like to /set new?")
	})

	log.Println("Handler for /dismiss - initialized.")
}

func initPingHandler(b telebot.Bot) {
	b.Handle("/ping", func(c telebot.Context) error {
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
	log.Println("Handler for /ping - initialized.")
}

func initHelpHandler(b telebot.Bot) {
	b.Handle("/help", func(c telebot.Context) error {
		return c.Send(descriptionMsg)
	})
	log.Println("Handler for /help - initialized.")
}

func initStartHandler(b telebot.Bot) {
	b.Handle("/start", func(c telebot.Context) error {
		var senderUsername = c.Sender().Username
		var helloMsg = fmt.Sprintf("Hello %s!", senderUsername)
		c.Send(helloMsg)
		return c.Send(descriptionMsg)
	})

	log.Println("Handler for /start - initialized.")
}

func initOnTextHandler(b telebot.Bot, err error) {
	b.Handle(telebot.OnText, func(c telebot.Context) error {
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

			log.Printf("Reminder for %s created!\n", senderUsername)
			return c.Send("Done. You can /ping for status.")
		}

		return err
	})

	log.Println("Text processor - initialized.")
}

func runReminderCronJob() {
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

	log.Println("Reminder Cron job - started.")
}
