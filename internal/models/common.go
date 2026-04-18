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
