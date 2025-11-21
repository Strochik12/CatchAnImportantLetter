package filter

import (
	"fmt"
	"log"

	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
)

type Engine struct {
	rules []*models.Rule
}

// NewEngine - создает новый движок правил
func NewEngine(r []*models.Rule) *Engine {
	return &Engine{
		rules: r,
	}
}

// Process - обрабатывает письмо через все правила
func (e *Engine) Process(email *models.Email) []*models.Alert {
	var alerts []*models.Alert = make([]*models.Alert, 0)

	for _, rule := range e.rules {
		if !rule.Enabled {
			continue
		}

		alert, err := e.evaluateRule(rule, email)
		if err != nil {
			log.Printf("Ошибка обработки правила: %v", err)
			continue
		}

		if alerts != nil {
			alerts = append(alerts, alert)
		}
	}

	return alerts
}

// evaluateRule - применяет одно правило к письму
func (e *Engine) evaluateRule(rule *models.Rule, email *models.Email) (*models.Alert, error) {
	score, reasons, err := e.calculateScore(rule, email)
	if err != nil {
		return nil, err
	}

	if score >= rule.MinScore {
		reasonText := fmt.Sprintf("Правило: %s. Баллы: %d/%d. Причины: %s",
			rule.Name, score, rule.MinScore, reasons)

		alert := models.NewAlert(email, rule, score, reasonText)
		return alert, nil
	}

	return nil, nil
}
