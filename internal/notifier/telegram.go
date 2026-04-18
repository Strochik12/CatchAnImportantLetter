package notifier

import (
	"fmt"
	"log"
	"strings"

	"github.com/Strochik12/CatchAnImportantLetter/internal/config"
	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramNotifier struct {
	BaseNotifier
	bot     *tgbotapi.BotAPI
	chatID  int64
	enabled bool
}

func NewTelegram(cfg *config.TelegramConfig) (*TelegramNotifier, error) {
	if cfg == nil || !cfg.Enabled || cfg.BotToken == "" || cfg.ChatID == 0 {
		return &TelegramNotifier{
			BaseNotifier: BaseNotifier{name: "telegram"},
			enabled:      false,
		}, nil
	}

	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания Telegram бота: %w", err)
	}

	notifier := &TelegramNotifier{
		BaseNotifier: BaseNotifier{name: "telegram"},
		bot:          bot,
		chatID:       cfg.ChatID,
		enabled:      true,
	}

	// Проверяем подключение
	if err := notifier.testConnection(); err != nil {
		return nil, fmt.Errorf("ошибка подключения к Telegram: %w", err)
	}

	log.Printf("Telegram нотификатор инициализирован для чата: %d", notifier.chatID)
	return notifier, nil
}

// Send отправляет уведомление в Telegram
func (t *TelegramNotifier) Send(alert *models.Alert) error {
	if !t.enabled {
		return fmt.Errorf("telegram нотификатор отключен")
	}

	message := t.formatMessage(alert)

	msg := tgbotapi.NewMessage(t.chatID, message)
	msg.ParseMode = "HTML"

	_, err := t.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("ошибка отправки в Telegram: %w", err)
	}

	log.Printf("Уведомление отправлено в Telegram: %s", alert.Rule.Name)
	return nil
}

// formatMessage форматирует сообщение для Telegram
func (t *TelegramNotifier) formatMessage(alert *models.Alert) string {
	var sb strings.Builder

	// Эмодзи в зависимости от уровня важности
	var emoji string
	switch alert.Level {
	case models.AlertCritical:
		emoji = "🔴"
	case models.AlertHigh:
		emoji = "🟠"
	case models.AlertMedium:
		emoji = "🟡"
	case models.AlertLow:
		emoji = "🔵"
	default:
		emoji = "🔔"
	}

	sb.WriteString(fmt.Sprintf("%s <b>%s</b>\n", emoji, escapeHTML(alert.Rule.Name)))
	sb.WriteString(fmt.Sprintf("<b>Тема:</b> %s\n", escapeHTML(alert.Email.Subject)))
	sb.WriteString(fmt.Sprintf("<b>От:</b> %s\n", escapeHTML(alert.Email.From)))

	sb.WriteString(fmt.Sprintf("<b>Время:</b> %s\n", alert.Email.Date.Format("15:04 02.01")))
	// sb.WriteString(fmt.Sprintf("<b>Причина:</b> %s\n", escapeHTML(alert.Reason)))

	return sb.String()
}

// testConnection проверяет подключение к Telegram
func (t *TelegramNotifier) testConnection() error {
	msg := tgbotapi.NewMessage(t.chatID, "Запущен и готов к работе!")
	_, err := t.bot.Send(msg)
	return err
}

// IsAvailable проверяет доступность нотификатора
func (t *TelegramNotifier) IsAvailable() bool {
	return t.enabled
}

// Enable/Disable для управления состоянием
func (t *TelegramNotifier) Enable()  { t.enabled = true }
func (t *TelegramNotifier) Disable() { t.enabled = false }

// escapeHTML экранирует HTML символы
func escapeHTML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)

	return replacer.Replace(text)
}
