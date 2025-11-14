package models

import (
	"crypto/rand"
	"fmt"
)

type ID string

type ConditionType string

const (
	ConditionFrom    ConditionType = "from"
	ConditionSubject ConditionType = "subject"
	ConditionBody    ConditionType = "body"
	ConditionHeader  ConditionType = "header"
)

type ActionType string

const (
	ActionNotifyTelegram ActionType = "telegram"
	ActionNotifySms      ActionType = "sms"
	ActionAutoClick      ActionType = "autoclick"
)

type Operator string

const (
	OperatorContains   Operator = "contains"
	OperatorEquals     Operator = "equals"
	OperatorStartsWith Operator = "startswith"
	OperatorEndsWith   Operator = "endswith"
	OperatorMatches    Operator = "matches"
)

// GenerateID - генерирует случайный ID
func GenerateID() ID {
	b := make([]byte, 16)
	rand.Read(b)
	return ID(fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]))
}

// ValidateRule проверяет правило на валидность
func (r *Rule) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("Rule name cannot be empty")
	}
	if len(r.Conditions) == 0 {
		return fmt.Errorf("Rule must have at least one condition")
	}
	if len(r.Actions) == 0 {
		return fmt.Errorf("Rule must have at least one action")
	}
	if r.MinScore < 0 || r.MinScore > 100 {
		return fmt.Errorf("min_score must be between 0 and 100")
	}
	return nil
}
