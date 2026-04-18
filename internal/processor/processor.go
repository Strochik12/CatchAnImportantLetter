package processor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Strochik12/CatchAnImportantLetter/internal/config"
	"github.com/Strochik12/CatchAnImportantLetter/internal/filter"
	"github.com/Strochik12/CatchAnImportantLetter/internal/mailwatcher"
	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
	"github.com/Strochik12/CatchAnImportantLetter/internal/notifier"
)

// Processor главный координатор системы
type Processor struct {
	config   *config.Config
	watcher  *mailwatcher.Watcher
	filter   *filter.Engine
	notifier *notifier.Manager
	stats    *Stats
}

type Stats struct {
	EmailsProcessed   int
	AlertsGenerated   int
	NotificationsSent int
	LastActivity      time.Time
	Errors            []error
}

// NewProcessor создаёт новый обработчик
func NewProcessor(cfg *config.Config) (*Processor, error) {
	watcher := mailwatcher.NewWatcher(cfg)

	filter := filter.NewEngine(cfg.Rules)

	notifier, err := notifier.NewManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания менеджера нотификаторов: %w", err)
	}

	return &Processor{
		config:   cfg,
		watcher:  watcher,
		filter:   filter,
		notifier: notifier,
		stats:    &Stats{LastActivity: time.Now()},
	}, nil
}

// Start запускает мониторинг почты
func (p *Processor) Start(ctx context.Context) error {
	log.Println("Запускаем HSE Email Alert System...")
	log.Printf("Сервер: %s", p.config.IMAP.Server)
	log.Printf("Пользователь: %s", p.config.IMAP.Username)
	log.Printf("Правил загружено: %d", len(p.config.Rules))
	log.Printf("Доступные нотификаторы: %v", p.notifier.GetAvailableNotifiers())

	// Запускаем мониторинг почты
	emailCh, errorCh := p.watcher.Watch(ctx)

	// Основной цикл обработки
	for {
		select {
		case email, ok := <-emailCh:
			if !ok {
				return nil
			}

			p.stats.LastActivity = time.Now()
			p.stats.EmailsProcessed += 1

			log.Printf("Новое письмо: %q", email.Subject)

			// Обрабатываем письмо
			if err := p.processEmail(email); err != nil {
				log.Printf("Ошибка обработки письма: %v", err)
				p.stats.Errors = append(p.stats.Errors, err)
			}

		case err, ok := <-errorCh:
			if !ok {
				return nil
			}
			log.Printf("Ошибка мониторинга: %v", err)
			p.stats.Errors = append(p.stats.Errors, err)

		case <-ctx.Done():
			log.Println("Останавливаем систему...")
			return nil
		}
	}
}

// processEmail обрабатывает одно письмо end-to-end
func (p *Processor) processEmail(email *models.Email) error {
	startTime := time.Now()

	// 1. Фильтруем через движок правил
	results := p.filter.Process(email)

	if len(results) == 0 {
		log.Printf("	Письмо не подошло ни под одно правило")
		return nil
	}

	log.Printf("	Сработавших правил: %d", len(results))
	p.stats.AlertsGenerated += len(results)

	// 2. Группируем действия по типу
	actionsByType := make(map[models.ActionType]bool)

	for _, result := range results {
		if result == nil {
			log.Printf("⚠️ Получен nil alert")
			continue
		}

		if result.Rule == nil {
			log.Printf("⚠️ Получен alert с nil Rule")
			continue
		}

		for _, action := range result.Rule.Actions {
			actionsByType[action] = true
		}
	}

	// 3. Выполняем действия
	var sentCount int
	var errors []error

	for actionType, _ := range actionsByType {
		if !p.notifier.HasNotifier(actionType) {
			log.Printf("	Нотификатор для %s, недоступен", actionType)
			continue
		}

		// Создаем алерт для отправки (берем первый из результатов)
		alert := results[0]

		if alert == nil {
			log.Printf("⚠️ Невозможно отправить alert, так как первый полученный alert - nil")
			break
		}

		// Отправляем уведомление
		if err := p.notifier.Send(actionType, alert); err != nil {
			log.Printf("	Ошибка отправки %s: %v", actionType, err)
			errors = append(errors, err)
		} else {
			sentCount++
			p.stats.NotificationsSent++
		}
	}

	processingTime := time.Since(startTime)
	log.Printf("	Результат: %d уведомлений отправлено за %v", sentCount, processingTime)

	if len(errors) > 0 {
		return fmt.Errorf("ошибки при отправке уведомлений: %v", errors)
	}

	return nil
}

// GetStats возвращает статистику работы
func (p *Processor) GetStats() *Stats {
	return p.stats
}

// PrintStats выводит статистику в консоль
func (p *Processor) PrintStats() {
	stats := p.stats
	fmt.Println("\nСтатистика работы:")
	fmt.Printf("	Обработано писем: %d\n", stats.EmailsProcessed)
	fmt.Printf("	Сгенерировано алертов: %d\n", stats.AlertsGenerated)
	fmt.Printf("	Отправлено уведомлений: %d\n", stats.NotificationsSent)
	fmt.Printf("	Последняя активность: %v\n", stats.LastActivity.Format("15:04:05"))

	if len(stats.Errors) > 0 {
		fmt.Printf("	Ошибок: %d\n", len(stats.Errors))
		for i, err := range stats.Errors {
			fmt.Printf("		%d. %v\n", i+1, err)
		}
	}
}
