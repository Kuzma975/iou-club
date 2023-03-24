package main

import (
	"log"
	"os"

	// "errors"
	// "database/sql"
	// _ "github.com/mattn/go-sqlite3"
	"kuzma975/iou-club/database"
	"kuzma975/iou-club/database/logging"
	"kuzma975/iou-club/handler"

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

func main() {
	logFile := logging.InitializeLogging(database.LogFile)
	defer func() {
		err := logFile.Close()
		logging.CheckErr(err)
		logging.Info.Printf("Logfile %s is closed", database.LogFile)
	}()

	db := database.InitializeDatabase()
	defer func() {
		err := db.Close()
		logging.CheckErr(err)
		logging.Info.Printf("Database %s is closed\n", database.DatabaseFile)
	}()

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
		log.Panic("Error occured during creating new bot client: ", err)
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
			handler.HandleMessage(start, *bot, db, test)
			defer func() {
				log.Printf("latest id is %d", latest.MessageID)
				log.Printf("first id %d", start.Message.MessageID)
				for i := start.Message.MessageID; i <= latest.MessageID; i++ {
					toDelete := tgbotapi.NewDeleteMessage(start.Message.Chat.ID, i)
					if resp, err := bot.Request(toDelete); err != nil {
						log.Printf("Cloud not delete message %s", err)
					} else {
						log.Printf("Response is %v", resp)
					}
				}
			}()
		}
	}
	for update := range updates {
		if update.Message != nil { // If we got a message
			latest = update.Message
			if handler.HandleMessage(update, *bot, db, false) {
				return
			}
		}
	}
}
