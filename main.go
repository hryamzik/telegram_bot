package main

import (
	"flag"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"gopkg.in/gin-gonic/gin.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Config struct {
	TelegramToken string `yaml:"telegram_token"`
}

var config_path = flag.String("c", "config.yaml", "Path to a config file")
var listen_addr = flag.String("l", ":9037", "Listen address")
var debug = flag.Bool("d", false, "debug")

var cfg = Config{}
var bot *tgbotapi.BotAPI

func main() {
	log.Println("Starting...")
	defer log.Println("Finished")
	flag.Parse()
	content, err := ioutil.ReadFile(*config_path)

	if err != nil {
		log.Fatalf("Problem reading configuration file: %v", err)
	}
	err = yaml.Unmarshal(content, &cfg)
	if err != nil {
		log.Fatalf("Error parsing configuration file: %v", err)
	}

	bot, err = tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	if *debug {
		bot.Debug = true
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	log.Printf("Authorised on account %s", bot.Self.UserName)

	go telegramBot(bot)

	router := gin.Default()
	router.POST("/:chatid", handler)
	router.Run(*listen_addr)
}

func handler(c *gin.Context) {
	if *debug {
		log.Println("Got POST request, looking for chat id")
	}
	chatid, err := strconv.ParseInt(c.Param("chatid"), 10, 64)
	if err != nil {
		log.Printf("Cat't parse chat id: %q\n", c.Param("chatid"))
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"err": fmt.Sprint(err),
		})
		return
	}
	if *debug {
		log.Printf("Should post to chat %d, looking for posted text\n", chatid)
	}

	msgtext, err := ioutil.ReadAll(c.Request.Body)

	if err != nil {
		log.Printf("Cat't get posted data\n")
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"err": fmt.Sprint(err),
		})
		return
	}

	if *debug {
		log.Printf("Got posted data, %d bites, text:\n", len(msgtext))
		log.Printf("=== %s ===\nChecking message format\n", msgtext)
	}

	mode := c.DefaultQuery("mode", "HTML")

	msg := tgbotapi.NewMessage(chatid, string(msgtext))

	if mode != "HTML" {
		msg.ParseMode = tgbotapi.ModeMarkdown
	} else {
		msg.ParseMode = tgbotapi.ModeHTML
	}

	if *debug {
		log.Printf("Message format is %s, sending message...\n", msg.ParseMode)
	}

	msg.DisableWebPagePreview = true

	sendmsg, err := bot.Send(msg)

	if *debug {
		log.Printf("Message sent\n")
	}

	if err == nil {
		c.String(http.StatusOK, "telegram msg sent.")
	} else {
		log.Printf("Error sending message: %s", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"err":     fmt.Sprint(err),
			"message": sendmsg,
			"srcmsg":  fmt.Sprint(msgtext),
		})
		msg := tgbotapi.NewMessage(chatid, "Error sending message, checkout logs")
		bot.Send(msg)
	}

}

func telegramBot(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	introduce := func(update tgbotapi.Update) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Chat id is '%d'", update.Message.Chat.ID))
		bot.Send(msg)
	}

	for update := range updates {
		if update.Message.NewChatMembers != nil && len(*update.Message.NewChatMembers) > 0 {
			for _, member := range *update.Message.NewChatMembers {
				if member.UserName == bot.Self.UserName && update.Message.Chat.Type == "group" {
					introduce(update)
				}
			}
		} else if update.Message != nil && update.Message.Text != "" {
			introduce(update)
		}
	}
}
