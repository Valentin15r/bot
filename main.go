package main

import (
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// URL твоего сайта с мини-приложением.
// ВАЖНО: должен быть HTTPS — Telegram не принимает http.
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

	// Кнопка, открывающая мини-приложение.
	// Её показываем в /start и по команде /app.
	webAppButton := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp(
				"🚀 Открыть приложение",
				tgbotapi.WebAppInfo{URL: webAppURL},
			),
		),
	)

	// Reply-кнопка (внизу экрана), тоже открывает WebApp.
	// Работает только в личных чатах.
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
					"Привет! 👋\nОткрой мини-приложение кнопкой ниже.")
				reply.ReplyMarkup = webAppButton
				reply.ReplyToMessageID = msg.MessageID
				send(bot, reply)

			case "app":
				reply := tgbotapi.NewMessage(msg.Chat.ID, "Жми 👇")
				reply.ReplyMarkup = menuButton
				send(bot, reply)

			case "help":
				send(bot, tgbotapi.NewMessage(msg.Chat.ID,
					"Команды:\n"+
						"/start — приветствие + кнопка WebApp\n"+
						"/app — показать reply-кнопку с приложением\n"+
						"/help — помощь\n"+
						"/time — время"))

			case "time":
				now := time.Now().Format("02.01.2006 15:04:05")
				send(bot, tgbotapi.NewMessage(msg.Chat.ID, "🕐 "+now))

			default:
				send(bot, tgbotapi.NewMessage(msg.Chat.ID, "Неизвестная команда. /help"))
			}
			continue
		}

		// Обработка данных, отправленных из WebApp через Telegram.WebApp.sendData()
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