package models

import (
	"fmt"
	"time"
)

type Alert struct {
	ID        ID         `json:"id"`
	Email     *Email     `json:"email"`
	Rule      *Rule      `json:"rule"`
	Score     int        `json:"score"`
	Level     AlertLevel `json:"level"`
	Reason    string     `json:"reason"`
	Message   string     `json:"message"`
	CreatedAt time.Time  `json:"created_at"`
	Processed bool       `json:"processed"`
}

type AlertLevel int

const (
	AlertLow      AlertLevel = 1
	AlertMedium   AlertLevel = 2
	AlertHigh     AlertLevel = 3
	AlertCritical AlertLevel = 4
)

// NewAlert создает новые Alert
func NewAlert(e *Email, r *Rule, score int, reason string) *Alert {
	alert := Alert{
		ID:        GenerateID(),
		Email:     e,
		Rule:      r,
		Score:     score,
		Reason:    reason,
		CreatedAt: time.Now(),
		Processed: false,
	}

	// Определяем уровень важности на основе баллов
	alert.calculateLevel()
	alert.generateMessage()

	return &alert
}

// calculateLevel определяет уровень важности
func (a *Alert) calculateLevel() {
	switch {
	case a.Score >= 90:
		a.Level = AlertCritical
	case a.Score >= 75:
		a.Level = AlertHigh
	case a.Score >= 65:
		a.Level = AlertMedium
	default:
		a.Level = AlertLow
	}
}

// generateMessage генерирует сообщение для уведомления
func (a *Alert) generateMessage() {
	levelNames := map[AlertLevel]string{
		AlertLow:      "🔔",
		AlertMedium:   "⚠️",
		AlertHigh:     "🚨",
		AlertCritical: "🔥",
	}

	a.Message = fmt.Sprintf("%s %s\nТема: %s\nОт: %s\nБалл: %d/%d\nПричина: %s",
		levelNames[a.Level],
		a.Rule.Name,
		a.Email.Subject,
		a.Email.From,
		a.Score,
		100,
		a.Reason,
	)
}

// MarkProcessed отмечает алерт как обработанный
func (a *Alert) MarkProcessed() {
	a.Processed = true
}
