package notifier

import (
	"fmt"
	"log"

	"github.com/Strochik12/CatchAnImportantLetter/internal/config"
	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
)

type Manager struct {
	notifiers map[models.ActionType]Notifier
}

// NewManager создает и настраивает все нотификаторы из конфига
func NewManager(cfg *config.Config) (*Manager, error) {
	manager := &Manager{
		notifiers: make(map[models.ActionType]Notifier),
	}

	// Инициализируем Telegram нотификатор
	if cfg.Notifiers.Telegram != nil && cfg.Notifiers.Telegram.Enabled {
		telegram, err := NewTelegram(cfg.Notifiers.Telegram)
		if err != nil {
			return nil, fmt.Errorf("ошибка инициализации Telegram:%w", err)
		}
		manager.Register(models.ActionNotifyTelegram, telegram)
	}

	// Здесь можно добавить инициализацию других нотификаторов
	// if cfg.Notifiers.SMS != nil && cfg.Notifiers.SMS.Enabled {
	//     sms, err := NewSMS(cfg.Notifiers.SMS)
	//     manager.Register(models.ActionNotifySMS, sms)
	// }

	log.Printf("Менеджер нотификаторов инициализирован. Доступно: %d", len(manager.notifiers))
	return manager, nil
}

func (m *Manager) Register(actionType models.ActionType, notifier Notifier) {
	if notifier.IsAvailable() {
		m.notifiers[actionType] = notifier
		log.Printf("Зарегистрирован нотификатор: %s", notifier.Name())
	}
}

// Send отправляет уведомление через соответствующий нотификатор
func (m *Manager) Send(actionType models.ActionType, alert *models.Alert) error {
	notifier, exists := m.notifiers[actionType]
	if !exists {
		return fmt.Errorf("нотификатор для действия %s не указан", actionType)
	}

	return notifier.Send(alert)
}

func (m *Manager) GetAvailableNotifiers() []string {
	var available []string
	for actionType, notifier := range m.notifiers {
		if notifier.IsAvailable() {
			available = append(available, string(actionType))
		}
	}
	return available
}

// HasNotifier проверяет наличие нотификатора для типа действия
func (m *Manager) HasNotifier(actionType models.ActionType) bool {
	notifier, exists := m.notifiers[actionType]
	return exists && notifier.IsAvailable()
}
