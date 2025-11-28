package main

import (
	"log"
	"time"

	"github.com/Strochik12/CatchAnImportantLetter/internal/config"
	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
	"github.com/Strochik12/CatchAnImportantLetter/internal/notifier"
)

func main() {
	// Загружаем конфиг
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	// Создаём менджер нотификаторов
	manager, err := notifier.NewManager(cfg)
	if err != nil {
		log.Fatalf("Ошибка создания менеджера нотификаторов: %v", err)
	}

	log.Printf("Доступные нотификаторы: %v", manager.GetAvailableNotifiers())

	// Тестовый алерт
	testAlert := &models.Alert{
		ID:   models.GenerateID(),
		Rule: &models.Rule{Name: "Тестовое правило"},
		Email: &models.Email{
			From:    "med@hse.ru",
			Subject: "Тест менеджера нотификаторов",
			Body:    "Проверка работы централизованной системы уведомлений",
			Date:    time.Now(),
		},
		Score: 85,
		Level: models.AlertHigh,
	}

	// Отправляем через менеджер
	if manager.HasNotifier(models.ActionNotifyTelegram) {
		log.Println("Отправляю тестовое уведомление через менеджер...")
		if err := manager.Send(models.ActionNotifyTelegram, testAlert); err != nil {
			log.Fatalf("Ошибка отправки: %v", err)
		}
		log.Println("Уведомление отправлено!")
	} else {
		log.Println("Telegram нотификатор недоступен")
	}
}
