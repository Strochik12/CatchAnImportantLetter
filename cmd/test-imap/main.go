package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Strochik12/CatchAnImportantLetter/internal/config"
	"github.com/Strochik12/CatchAnImportantLetter/internal/mailwatcher"
)

func main() {
	// Загружаем конфиг
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	// Создаем Watcher
	watcher := mailwatcher.NewWatcher(cfg)

	// Настраиваем gracefull shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Println("Запускается тест IMAP клиента...")
	log.Printf("Сервер: %s", cfg.IMAP.Server)
	log.Printf("Пользователь: %s", cfg.IMAP.Username)

	// Запускаем мониторинг
	emailCh, errorCh := watcher.Watch(ctx)

	// Обрабатываем результаты
	for {
		select {
		case email, ok := <-emailCh:
			if !ok {
				return
			}
			log.Printf("Новое письмо от: %s", email.From)
			log.Printf("Дата: %s", email.Date)
			if email.Body != "" {
				bodyPreview := email.Body
				if len(bodyPreview) > 100 {
					bodyPreview = bodyPreview[:100] + "..."
				}
				log.Printf("Текст: %s", bodyPreview)
			}

		case err, ok := <-errorCh:
			if !ok {
				return
			}
			log.Printf("Ошибка: %v", err)

		case <-ctx.Done():
			log.Println("Завершаем работу...")
			return
		}
	}
}
