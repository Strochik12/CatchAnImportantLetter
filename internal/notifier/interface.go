package notifier

import "github.com/Strochik12/CatchAnImportantLetter/internal/models"

// Notifier интерфейс нотификатора для всех типов уведомлений
type Notifier interface {
	Send(alert *models.Alert) error
	Name() string
	IsAvailable() bool
}

//BaseNotifier базовая структура для всех нотификаторов
type BaseNotifier struct {
	name string
}

func (b *BaseNotifier) Name() string {
	return b.name
}
