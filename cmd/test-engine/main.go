package main

import (
	"fmt"
	"log"

	"github.com/Strochik12/CatchAnImportantLetter/internal/config"
	"github.com/Strochik12/CatchAnImportantLetter/internal/filter"
	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
)

func main() {
	// Загружаем конфиг
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("ошибка загрузки конфига: %v", err)
	}

	// Создаем движок правил
	engine := filter.NewEngine(cfg.Rules)

	// Тестовое письмо
	testEmail := models.Email{
		Subject: "Запись на медосмотр для сотрудников",
		From:    "med@hse.ru",
		Body:    "Уважаемые сотрудники, открыта запись на ежегодный медосмотр!",
	}

	fmt.Printf("Обрабатываем письмо:\n")
	fmt.Printf("  От: %s\n", testEmail.From)
	fmt.Printf("  Тема: %s\n", testEmail.Subject)

	// Обрабатываем письмо
	alerts := engine.Process(&testEmail)

	fmt.Printf("\nРезультаты:\n")
	if len(alerts) == 0 {
		fmt.Println("	Не создано алертов")
	} else {
		for i, alert := range alerts {
			fmt.Printf("  %d. %s\n", i+1, alert.Message)
			fmt.Printf("     Правило: %s, Баллы: %d\n", alert.Rule.Name, alert.Score)
		}
	}

	// Показываем все правила
	fmt.Printf("\nЗагруженные правила (%d):\n", len(cfg.Rules))
	for _, rule := range cfg.Rules {
		status := "enabled "
		if !rule.Enabled {
			status = "disabled "
		}
		fmt.Printf("  %s %s (мин. балл: %d, приоритет: %d)\n",
			status, rule.Name, rule.MinScore, rule.Priority)
	}
}
