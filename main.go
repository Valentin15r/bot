package main

import (
	"log"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Токен берём из переменной окружения BOT_TOKEN
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatal("Не задан BOT_TOKEN. Пример: export BOT_TOKEN=123456:ABC-DEF...")
	}

	// Создаём бота
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatalf("Ошибка создания бота: %v", err)
	}

	bot.Debug = false
	log.Printf("Бот запущен: @%s", bot.Self.UserName)

	// Получаем обновления
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		// Игнорируем обновления без сообщений
		if update.Message == nil {
			continue
		}

		msg := update.Message
		log.Printf("[%s] %s", msg.From.UserName, msg.Text)

		// Обработка команд
		if msg.IsCommand() {
			switch msg.Command() {
			case "start":
				reply := tgbotapi.NewMessage(msg.Chat.ID,
					"Привет! 👋\nЯ простой бот на Go.\n\n"+
						"Команды:\n"+
						"/start — приветствие\n"+
						"/help — помощь\n"+
						"/echo <текст> — повторить текст\n"+
						"/time — текущее время\n\n"+
						"Или просто напиши мне что-нибудь.")
				reply.ReplyToMessageID = msg.MessageID
				send(bot, reply)

			case "help":
				reply := tgbotapi.NewMessage(msg.Chat.ID,
					"Доступные команды:\n"+
						"/start — приветствие\n"+
						"/help — помощь\n"+
						"/echo <текст> — повторить текст\n"+
						"/time — текущее время")
				send(bot, reply)

			case "echo":
				args := msg.CommandArguments()
				if strings.TrimSpace(args) == "" {
					send(bot, tgbotapi.NewMessage(msg.Chat.ID, "Использование: /echo <текст>"))
				} else {
					send(bot, tgbotapi.NewMessage(msg.Chat.ID, args))
				}

			case "time":
				now := time.Now().Format("02.01.2006 15:04:05")
				send(bot, tgbotapi.NewMessage(msg.Chat.ID, "🕐 "+now))

			default:
				send(bot, tgbotapi.NewMessage(msg.Chat.ID, "Неизвестная команда. Попробуй /help"))
			}
			continue
		}

		// Эхо на обычный текст
		if msg.Text != "" {
			reply := tgbotapi.NewMessage(msg.Chat.ID, "Ты написал: "+msg.Text)
			reply.ReplyToMessageID = msg.MessageID
			send(bot, reply)
		}
	}
}

// send — обёртка для отправки с логированием ошибок
func send(bot *tgbotapi.BotAPI, msg tgbotapi.MessageConfig) {
	msg.ParseMode = "" // можно "Markdown" или "HTML", если нужно
	if _, err := bot.Send(msg); err != nil {
		log.Printf("Ошибка отправки: %v", err)
	}
}