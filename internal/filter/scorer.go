package filter

import (
	"fmt"
	"strings"

	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
)

// calculateScore - считает баллы полученные письмом за правило
func (e *Engine) calculateScore(rule *models.Rule, email *models.Email) (int, string, error) {
	var score int = 0
	var reasons []string = make([]string, 0)

	for _, condition := range rule.Conditions {
		check, reason, err := e.evaluateCondition(condition, email)
		if err != nil {
			return 0, "", err
		}

		if check {
			score += condition.Weight
			reasons = append(reasons, reason)
		}
	}

	return score, strings.Join(reasons, ", "), nil
}

// evaluateCondition - проверяет выполняет ли письмо условие
func (e *Engine) evaluateCondition(cond models.Condition, email *models.Email) (bool, string, error) {
	var value string
	var fieldName string
	switch cond.Type {
	case models.ConditionBody:
		value = email.Body
		fieldName = "Тело"

	case models.ConditionFrom:
		value = email.From
		fieldName = "Отправитель"

	case models.ConditionHeader:
		value = email.Headers[cond.Field]
		fieldName = fmt.Sprintf("Заголовок %s", cond.Field)

	case models.ConditionSubject:
		value = email.Subject
		fieldName = "Тема"

	default:
		return false, "", fmt.Errorf("неизвестный тип условия: %s", cond.Operator)
	}

	check, err := e.checkCondition(cond, value)
	if err != nil {
		return false, "", err
	}

	if check {
		return true, fmt.Sprintf("%s содержит %s", fieldName, cond.Value), nil
	}

	return false, "", nil
}
