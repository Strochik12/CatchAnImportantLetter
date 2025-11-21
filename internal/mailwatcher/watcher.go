package mailwatcher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Strochik12/CatchAnImportantLetter/internal/config"
	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
)

type Watcher = Client

// NewWatcher создаёт новый Watcher
func NewWatcher(cfg *config.Config) *Watcher {
	return NewIMAPClient(cfg)
}

// Watch запускает мониторинг почты
func (w *Watcher) Watch(ctx context.Context) (<-chan *models.Email, <-chan error) {
	emailCh := make(chan *models.Email)
	errorCh := make(chan error)

	go func() {
		defer close(emailCh)
		defer close(errorCh)

		// Подключаемся к почте
		if err := w.Connect(); err != nil {
			errorCh <- fmt.Errorf("ошибка подключения: %w", err)
			return
		}
		defer w.Close()

		log.Printf("Мониторинг почты запущен (интервал: %v)", w.config.Monitoring.GetCheckInterval())

		ticker := time.NewTicker(w.config.Monitoring.GetCheckInterval())
		defer ticker.Stop()

		// Первый просмотр сразу при запуске, чтобы не ждать.
		emails, err := w.GetNewEmails()
		if err != nil {
			errorCh <- fmt.Errorf("ошибка проверки почты: %w", err)
		} else {
			for _, email := range emails {
				emailCh <- email
			}
		}

		for {
			select {
			case <-ticker.C:
				emails, err := w.GetNewEmails()
				if err != nil {
					errorCh <- fmt.Errorf("ошибка проверки почты: %w", err)
					continue
				}

				for _, email := range emails {
					emailCh <- email
				}
			case <-ctx.Done():
				log.Println("Мониторинг почты остановлен")
				return
			}
		}
	}()

	return emailCh, errorCh
}
