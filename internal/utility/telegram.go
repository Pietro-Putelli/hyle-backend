package utility

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pietro-putelli/feynman-backend/config"
	"github.com/pietro-putelli/feynman-backend/internal/domain"
)

func sendMessage(message string) {
	cfg, err := config.NewConfig()
	if err != nil {
		return
	}

	bot, err := tgbotapi.NewBotAPI(cfg.Telegram.ApiToken)
	if err != nil {
		fmt.Println(err)
	}

	bot.Send(tgbotapi.NewMessage(-4264723912, message))
}

func TelegramSendNewUser(user *domain.User) {
	sendMessage("New user: " + user.GivenName + " " + user.FamilyName)
}
