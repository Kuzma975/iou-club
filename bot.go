package main

import (
	"log"
	"os"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

type config struct {
	Telegram struct {
		Token string `yaml:"token"`
		Debug bool   `yaml:"debug"`
		Test  bool   `yaml:"test"`
	} `yaml:"telegram"`
}

func handleMessage(update tgbotapi.Update, bot tgbotapi.BotAPI, isTest bool) bool {
	var msg tgbotapi.MessageConfig
	// var msgEnt tgbotapi.MessageEntity
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	if update.Message.Text == "exit" && !isTest {
		return true
	}
	if update.Message.ForwardFrom != nil {
		log.Printf("Account id is: %s", strconv.FormatInt(update.Message.ForwardFrom.ID, 10))
		// msgEnt = tgbotapi.MessageEntity{Type: "text_mention", User: update.Message.ForwardFrom}
	} else {
		log.Printf("Is private account")
	}
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "") // Hello, [user_name](tg://user?id="+strconv.FormatInt(update.Message.ForwardFrom.ID, 10)+")")
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "list":
			msg.Text = "/custom command"
		case "custom":
			msg.Text = "test is passed"
		default:
			msg.Text = "command not implemented"
		}
	}
	msg.ReplyToMessageID = update.Message.MessageID
	// msg.ParseMode = "MarkdownV2"
	// msg.Entities = []tgbotapi.MessageEntity{msgEnt}

	bot.Send(msg)
	return false
}

func main() {
	configurationFileName := "config.yaml"
	file, err := os.Open(configurationFileName)
	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	var conf config
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Panic(err)
	}
	// log.Printf("config file is: %s", conf)
	bot, err := tgbotapi.NewBotAPI(conf.Telegram.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = conf.Telegram.Debug

	log.Printf("Authorized on account %s", bot.Self.UserName)
	test := conf.Telegram.Test
	offset := func() int {
		if test {
			return -1
		} else {
			return 0
		}
	}()
	u := tgbotapi.NewUpdate(offset)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	var latest *tgbotapi.Message
	if test {
		start := <-updates
		if start.Message != nil {
			handleMessage(start, *bot, test)
			defer func() {
				log.Printf("latest id is %d", latest.MessageID)
				log.Printf("first id %d", start.Message.MessageID)
				for i := start.Message.MessageID; i <= latest.MessageID; i++ {
					toDelete := tgbotapi.NewDeleteMessage(start.Message.Chat.ID, i)
					if resp, err := bot.Request(toDelete); err != nil {
						log.Printf("Cloud not delete message %s", err)
					} else {
						log.Printf("Response is %s", resp)
					}
				}
			}()
		}
	}
	for update := range updates {
		if update.Message != nil { // If we got a message
			latest = update.Message
			if handleMessage(update, *bot, false) {
				return
			}
		}
	}
}
