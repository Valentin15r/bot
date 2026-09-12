package main

import (
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Впиши адрес своего сайта. Обязательно HTTPS.
const webAppURL = "https://example.com/app"

func main() {
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("Не задан BOT_TOKEN")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}
	log.Printf("Бот запущен: @%s", bot.Self.UserName)

	// Inline-кнопка, открывающая WebApp
	webAppButton := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp(
				"🚀 Открыть приложение",
				tgbotapi.WebAppInfo{URL: webAppURL},
			),
		),
	)

	// Reply-кнопка внизу экрана
	menuButton := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonWebApp(
				"📱 Приложение",
				tgbotapi.WebAppInfo{URL: webAppURL},
			),
		),
	)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		msg := update.Message

		if msg.IsCommand() {
			switch msg.Command() {
			case "start":
				reply := tgbotapi.NewMessage(msg.Chat.ID,
					"Привет! 👋 Открой мини-приложение:")
				reply.ReplyMarkup = webAppButton
				reply.ReplyToMessageID = msg.MessageID
				send(bot, reply)

			case "app":
				reply := tgbotapi.NewMessage(msg.Chat.ID, "Жми 👇")
				reply.ReplyMarkup = menuButton
				send(bot, reply)

			case "help":
				send(bot, tgbotapi.NewMessage(msg.Chat.ID,
					"Команды:\n/start — приветствие\n/app — кнопка приложения\n/time — время"))

			case "time":
				now := time.Now().Format("02.01.2006 15:04:05")
				send(bot, tgbotapi.NewMessage(msg.Chat.ID, "🕐 "+now))

			default:
				send(bot, tgbotapi.NewMessage(msg.Chat.ID, "Неизвестная команда. /help"))
			}
			continue
		}

		// Данные из WebApp (sendData)
		if msg.WebAppData != nil {
			log.Printf("WebAppData от %d: %s", msg.From.ID, msg.WebAppData.Data)
			send(bot, tgbotapi.NewMessage(msg.Chat.ID,
				"Получил из приложения: "+msg.WebAppData.Data))
			continue
		}

		if msg.Text != "" {
			reply := tgbotapi.NewMessage(msg.Chat.ID, "Ты написал: "+msg.Text)
			reply.ReplyToMessageID = msg.MessageID
			send(bot, reply)
		}
	}
}

func send(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки: %v", err)
	}
}