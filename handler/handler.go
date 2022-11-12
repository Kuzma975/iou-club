package handler

import (
	"database/sql"
	"kuzma975/iou-club/database"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func HandleMessage(update tgbotapi.Update, bot tgbotapi.BotAPI, db *sql.DB, isTest bool) bool {
	var msg tgbotapi.MessageConfig
	// var msgEnt tgbotapi.MessageEntity
	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Reply: "+update.Message.Text) // Hello, [user_name](tg://user?id="+strconv.FormatInt(update.Message.ForwardFrom.ID, 10)+")")
	log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	log.Printf("Income Text messge: %s", update.Message.Text)
	log.Printf("Income messge object: %+v", update.Message)
	if update.Message.From != nil {
		log.Printf("======= FROM ========")
		log.Printf("Income messge from: %+v", update.Message.From)
		log.Printf("Income messge from ID: %+v", update.Message.From.ID)
		log.Printf("Income messge from user name: %+v", update.Message.From.UserName)
		log.Printf("Income messge from first name: %+v", update.Message.From.FirstName)
		log.Printf("Income messge from last name: %+v", update.Message.From.LastName)
		log.Printf("Income messge from LanguageCode: %+v", update.Message.From.LanguageCode)

		log.Printf("+++++++++++++++")
	}
	if update.Message.ForwardFrom != nil {
		log.Printf("=======FORWARD FROM ========")
		log.Printf("Income messge ForwardFrom user name: %+v", update.Message.ForwardFrom.UserName)
		log.Printf("Income messge ForwardFrom first name: %+v", update.Message.ForwardFrom.LastName)
		log.Printf("Income messge ForwardFrom last name: %+v", update.Message.ForwardFrom.FirstName)
		log.Printf("+++++++++++++++")
	}
	log.Printf("Income message forward sender name: %+v", update.Message.ForwardSenderName)
	log.Printf("Income messge chat: %+v", update.Message.Chat)
	if update.Message.Entities != nil {
		log.Printf("Entities: %+v", update.Message.Entities)
	} else {
		log.Printf("Entity not found")
	}
	if strings.ToLower(update.Message.Text) == "exit" && !isTest {
		msg.Text = "Ok"
		bot.Send(msg)
		return true
	}
	if update.Message.ForwardFrom != nil {
		log.Printf("Account id is: %s", strconv.FormatInt(update.Message.ForwardFrom.ID, 10))
		// msgEnt = tgbotapi.MessageEntity{Type: "custom_emoji", CustomEmojiId: "5373141891321699086"}
	} else {
		log.Printf("Is private account")
	}
	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "list":
			msg.Text = "/custom command"
		case "custom":
			msg.Text = "test is passed"
		case "start":
			if update.Message.From != nil {
				log.Printf("======= FROM ========")
				log.Printf("Income messge from: %+v", update.Message.From)
				log.Printf("Income messge from ID: %+v", update.Message.From.ID)
				log.Printf("Income messge from user name: %+v", update.Message.From.UserName)
				log.Printf("Income messge from first name: %+v", update.Message.From.FirstName)
				log.Printf("Income messge from last name: %+v", update.Message.From.LastName)
				log.Printf("Income messge from LanguageCode: %+v", update.Message.From.LanguageCode)

				log.Printf("+++++++++++++++")
				affectedRow := database.CreateUser(
					db,
					update.Message.From.ID,
					update.Message.From.UserName,
					update.Message.From.FirstName,
					update.Message.From.LastName,
					update.Message.Date,
				) // check if create user if not exists
				log.Printf("Affected row is: %+v", affectedRow)
				if update.Message.Chat.Type == "group" {
					// initiate group
					// add user to the group
					// add mentioned user to the group
				}
			}
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
