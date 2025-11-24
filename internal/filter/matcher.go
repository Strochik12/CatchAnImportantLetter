package filter

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Strochik12/CatchAnImportantLetter/internal/models"
)

// checkCondition - проверяет строку на соответствие условию
func (e *Engine) checkCondition(cond models.Condition, value string) (bool, error) {
	if value == "" {
		return false, nil
	}

	value = strings.ToLower(value)

	switch cond.Operator {
	case models.OperatorContains:
		return strings.Contains(value, cond.Value), nil

	case models.OperatorEquals:
		// Равность нечувствительно к регистру
		return strings.EqualFold(value, cond.Value), nil

	case models.OperatorStartsWith:
		return strings.HasPrefix(value, cond.Value), nil

	case models.OperatorEndsWith:
		return strings.HasSuffix(value, cond.Value), nil

	case models.OperatorMatches:
		re, err := regexp.Compile(cond.Value)
		if err != nil {
			return false, fmt.Errorf("неверное регулярное выражение: %w", err)
		}
		return re.MatchString(value), nil

	default:
		return false, fmt.Errorf("неизвестный оператор условия: %s", cond.Operator)
	}
}
